// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package dumpling

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"time"

	"github.com/pingcap/dm/dm/config"
	"github.com/pingcap/dm/dm/pb"
	"github.com/pingcap/dm/dm/unit"
	"github.com/pingcap/dm/pkg/log"
	"github.com/pingcap/dm/pkg/terror"
	"github.com/pingcap/dm/pkg/utils"

	"github.com/pingcap/dumpling/v4/export"
	"github.com/pingcap/errors"
	"github.com/pingcap/failpoint"
	filter "github.com/pingcap/tidb-tools/pkg/table-filter"
	"github.com/siddontang/go/sync2"
	"go.uber.org/zap"
)

// Dumpling dumps full data from a MySQL-compatible database
type Dumpling struct {
	cfg *config.SubTaskConfig

	logger log.Logger

	dumpConfig *export.Config
	closed     sync2.AtomicBool
}

// NewDumpling creates a new Dumpling
func NewDumpling(cfg *config.SubTaskConfig) *Dumpling {
	m := &Dumpling{
		cfg:    cfg,
		logger: log.With(zap.String("task", cfg.Name), zap.String("unit", "dump")),
	}
	return m
}

// Init implements Unit.Init
func (m *Dumpling) Init(ctx context.Context) error {
	var err error
	m.dumpConfig, err = m.constructArgs()
	m.detectSQLMode()
	return err
}

// Process implements Unit.Process
func (m *Dumpling) Process(ctx context.Context, pr chan pb.ProcessResult) {
	dumplingExitWithErrorCounter.WithLabelValues(m.cfg.Name, m.cfg.SourceID).Add(0)

	failpoint.Inject("dumpUnitProcessWithError", func(val failpoint.Value) {
		m.logger.Info("dump unit runs with injected error", zap.String("failpoint", "dumpUnitProcessWithError"), zap.Reflect("error", val))
		msg, ok := val.(string)
		if !ok {
			msg = "unknown process error"
		}
		pr <- pb.ProcessResult{
			IsCanceled: false,
			Errors:     []*pb.ProcessError{unit.NewProcessError(errors.New(msg))},
		}
		failpoint.Return()
	})

	begin := time.Now()
	errs := make([]*pb.ProcessError, 0, 1)

	failpoint.Inject("dumpUnitProcessForever", func() {
		m.logger.Info("dump unit runs forever", zap.String("failpoint", "dumpUnitProcessForever"))
		<-ctx.Done()
		failpoint.Return()
	})

	// NOTE: remove output dir before start dumping
	// every time re-dump, loader should re-prepare
	err := os.RemoveAll(m.cfg.Dir)
	if err != nil {
		m.logger.Error("fail to remove output directory", zap.String("directory", m.cfg.Dir), log.ShortError(err))
		errs = append(errs, unit.NewProcessError(terror.ErrDumpUnitRuntime.Delegate(err, "fail to remove output directory: "+m.cfg.Dir)))
		pr <- pb.ProcessResult{
			IsCanceled: false,
			Errors:     errs,
		}
		return
	}

	newCtx, cancel := context.WithCancel(ctx)
	err = export.Dump(newCtx, m.dumpConfig)
	cancel()

	if err != nil {
		if utils.IsContextCanceledError(err) {
			m.logger.Info("filter out error caused by user cancel")
		} else {
			dumplingExitWithErrorCounter.WithLabelValues(m.cfg.Name, m.cfg.SourceID).Inc()
			errs = append(errs, unit.NewProcessError(terror.ErrDumpUnitRuntime.Delegate(err, "")))
		}
	}

	isCanceled := false
	select {
	case <-ctx.Done():
		isCanceled = true
	default:
	}

	if len(errs) == 0 {
		m.logger.Info("dump data finished", zap.Duration("cost time", time.Since(begin)))
	} else {
		m.logger.Error("dump data exits with error", zap.Duration("cost time", time.Since(begin)),
			zap.String("error", utils.JoinProcessErrors(errs)))
	}

	pr <- pb.ProcessResult{
		IsCanceled: isCanceled,
		Errors:     errs,
	}
}

// Close implements Unit.Close
func (m *Dumpling) Close() {
	if m.closed.Get() {
		return
	}

	m.removeLabelValuesWithTaskInMetrics(m.cfg.Name)
	// do nothing, external will cancel the command (if running)
	m.closed.Set(true)
}

// Pause implements Unit.Pause
func (m *Dumpling) Pause() {
	if m.closed.Get() {
		m.logger.Warn("try to pause, but already closed")
		return
	}
	// do nothing, external will cancel the command (if running)
}

// Resume implements Unit.Resume
func (m *Dumpling) Resume(ctx context.Context, pr chan pb.ProcessResult) {
	if m.closed.Get() {
		m.logger.Warn("try to resume, but already closed")
		return
	}
	// just call Process
	m.Process(ctx, pr)
}

// Update implements Unit.Update
func (m *Dumpling) Update(cfg *config.SubTaskConfig) error {
	// not support update configuration now
	return nil
}

// Status implements Unit.Status
func (m *Dumpling) Status() interface{} {
	// NOTE: try to add some status, like dumped file count
	return &pb.DumpStatus{}
}

// Error implements Unit.Error
func (m *Dumpling) Error() interface{} {
	return &pb.DumpError{}
}

// Type implements Unit.Type
func (m *Dumpling) Type() pb.UnitType {
	return pb.UnitType_Dump
}

// IsFreshTask implements Unit.IsFreshTask
func (m *Dumpling) IsFreshTask(ctx context.Context) (bool, error) {
	return true, nil
}

// constructArgs constructs arguments for exec.Command
func (m *Dumpling) constructArgs() (*export.Config, error) {
	cfg := m.cfg
	db := cfg.From

	dumpConfig := export.DefaultConfig()
	// ret is used to record the unsupported arguments for dumpling
	var ret []string

	// block status addr because we already have it in DM, and if we enable it, may we need more ports for the process.
	dumpConfig.StatusAddr = ""

	dumpConfig.Host = db.Host
	dumpConfig.Port = db.Port
	dumpConfig.User = db.User
	dumpConfig.Password = db.Password
	dumpConfig.OutputDirPath = cfg.Dir // use LoaderConfig.Dir as output dir
	tableFilter, err := filter.ParseMySQLReplicationRules(cfg.BAList)
	if err != nil {
		return nil, err
	}
	dumpConfig.TableFilter = tableFilter
	dumpConfig.EscapeBackslash = true
	dumpConfig.CompleteInsert = true // always keep column name in `INSERT INTO` statements.
	dumpConfig.Logger = m.logger.Logger

	if cfg.Threads > 0 {
		dumpConfig.Threads = cfg.Threads
	}
	if cfg.ChunkFilesize != "" {
		dumpConfig.FileSize, err = parseFileSize(cfg.ChunkFilesize)
		if err != nil {
			m.logger.Warn("parsed some unsupported arguments", zap.Error(err))
			return nil, err
		}
	}
	if cfg.StatementSize > 0 {
		dumpConfig.StatementSize = cfg.StatementSize
	}
	if cfg.Rows > 0 {
		dumpConfig.Rows = cfg.Rows
	}
	if len(cfg.Where) > 0 {
		dumpConfig.Where = cfg.Where
	}

	if cfg.SkipTzUTC {
		// TODO: support skip-tz-utc
		ret = append(ret, "--skip-tz-utc")
	}

	if db.Security != nil {
		dumpConfig.Security.CAPath = db.Security.SSLCA
		dumpConfig.Security.CertPath = db.Security.SSLCert
		dumpConfig.Security.KeyPath = db.Security.SSLKey
	}

	extraArgs := strings.Fields(cfg.ExtraArgs)
	if len(extraArgs) > 0 {
		err := parseExtraArgs(&m.logger, dumpConfig, ParseArgLikeBash(extraArgs))
		if err != nil {
			m.logger.Warn("parsed some unsupported arguments", zap.Error(err))
			ret = append(ret, extraArgs...)
		}
	}

	// record exit position when consistency is none, to support scenarios like Aurora upstream
	if dumpConfig.Consistency == "none" {
		m.logger.Info("check consistency is none")
	}

	m.logger.Info("create dumpling", zap.Stringer("config", dumpConfig))
	if len(ret) > 0 {
		m.logger.Warn("meeting some unsupported arguments", zap.Strings("argument", ret))
	}

	if !cfg.CaseSensitive {
		dumpConfig.TableFilter = filter.CaseInsensitive(dumpConfig.TableFilter)
	}

	return dumpConfig, nil
}

// detectSQLMode tries to detect SQL mode from upstream. If success, write it to LoaderConfig
func (m *Dumpling) detectSQLMode() {
	db, err := sql.Open("mysql", m.dumpConfig.GetDSN(""))
	if err != nil {
		return
	}
	defer db.Close()

	sqlMode, err := utils.GetGlobalVariable(db, "sql_mode")
	if err != nil {
		return
	}
	m.logger.Info("found upstream SQL mode", zap.String("SQL mode", sqlMode))
	m.cfg.LoaderConfig.SQLMode = sqlMode
}
