module github.com/pingcap/dm

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e
	github.com/coreos/go-systemd v0.0.0-20190719114852-fd7a80b32e1f // indirect
	github.com/docker/go-units v0.4.0
	github.com/dustin/go-humanize v1.0.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/gateway v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/groupcache v0.0.0-20191227052852-215e87163ea7 // indirect
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.4
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.14.3
	github.com/kami-zh/go-capturer v0.0.0-20171211120116-e492ea43421d
	github.com/klauspost/cpuid v1.2.1 // indirect
	github.com/lance6716/retool v1.3.8-0.20200806070832-3469f70b2afe
	github.com/montanaflynn/stats v0.5.0 // indirect
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/pingcap/check v0.0.0-20200212061837-5e12011dc712
	github.com/pingcap/dumpling v0.0.0-20200822215635-6f0a11c00c88
	github.com/pingcap/errcode v0.3.0 // indirect
	github.com/pingcap/errors v0.11.5-0.20200820035142-66eb5bf1d1cd
	github.com/pingcap/failpoint v0.0.0-20200702092429-9f69995143ce
	github.com/pingcap/log v0.0.0-20200828042413-fce0951f1463
	github.com/pingcap/parser v0.0.0-20200427031542-879c7bd4f27d
	github.com/pingcap/tidb v1.1.0-beta.0.20210831131150-57f701820764
	github.com/pingcap/tidb-tools v4.0.6-0.20200819060105-3ac93f6b99d4+incompatible
	github.com/pingcap/tipb v0.0.0-20200618092958-4fad48b4c8c3 // indirect
	github.com/prometheus/client_golang v1.5.1
	github.com/rakyll/statik v0.1.6
	github.com/satori/go.uuid v1.2.0
	github.com/shopspring/decimal v0.0.0-20191125035519-b054a8dfd10d // indirect
	github.com/siddontang/go v0.0.0-20180604090527-bdc77568d726
	github.com/siddontang/go-log v0.0.0-20190221022429-1e957dd83bed // indirect
	github.com/siddontang/go-mysql v0.0.0-20200222075837-12e89848f047
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/soheilhy/cmux v0.1.4
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/syndtr/goleveldb v1.0.1-0.20190625010220-02440ea7a285
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/uber/jaeger-client-go v2.22.1+incompatible // indirect
	github.com/uber/jaeger-lib v2.2.0+incompatible // indirect
	github.com/unrolled/render v1.0.1
	go.etcd.io/bbolt v1.3.5 // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20191023171146-3cf2f69b5738
	go.uber.org/zap v1.16.0
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200824131525-c12d262b63d8
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac
	golang.org/x/tools v0.0.0-20200823205832-c024452afbcd // indirect
	google.golang.org/genproto v0.0.0-20191230161307-f3c370f40bfb
	google.golang.org/grpc v1.26.0
	gopkg.in/yaml.v2 v2.3.0
	honnef.co/go/tools v0.0.1-2020.1.5 // indirect
	sourcegraph.com/sourcegraph/appdash v0.0.0-20190731080439-ebfcffb1b5c0 // indirect
)

go 1.13

replace github.com/pingcap/kvproto => github.com/pingcap/kvproto v0.0.0-20200420075417-e0c6e8842f22

replace github.com/pingcap/pd/v4 => github.com/nolouch/pd/v4 v4.0.0-20210831114947-686590ed34cd
