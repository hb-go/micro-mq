package auth

import (
	"flag"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/hb-go/micro-mq/pkg/log"

	"github.com/hb-go/micro-mq/auth/handler"
	"github.com/hb-go/micro-mq/auth/proto"
	"github.com/hb-go/micro-mq/pkg/util/conv"
)

var (
	cmdHelp  = flag.Bool("h", false, "帮助")
	addr     = flag.String("addr", "localhost:8972", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
)

func init() {
	flag.Parse()
}

func main() {
	if *cmdHelp {
		flag.PrintDefaults()
		return
	}

	s := server.NewServer()
	addRegistryPlugin(s)

	s.RegisterName(proto.SRV_Auth.String(), new(handler.Auth), "")
	s.Serve("tcp", *addr)
}

func addRegistryPlugin(s *server.Server) {

	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       conv.ProtoEnumsToRpcxBasePath(proto.BASE_PATH_name),
		Metrics:        metrics.NewRegistry(),
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
