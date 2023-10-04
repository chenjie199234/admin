package dao

import (
	"crypto/tls"

	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/model"

	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
)

// var ExampleCGrpcApi example.ExampleCGrpcClient
// var ExampleCrpcApi example.ExampleCrpcClient
// var ExampleWebApi  example.ExampleWebClient
var DingTalkWebClient *web.WebClient
var FeiShuWebClient *web.WebClient

// NewApi create all dependent service's api we need in this program
func NewApi() error {
	var e error
	_ = e //avoid unuse

	//init dns discover for example server
	//exampleDnsDiscover, e := discover.NewDNSDiscover("exampleproject", "examplegroup", "examplename", "dnshost", time.Second*10, 9000, 10000, 8000)
	//if e != nil {
	//	return e
	//}
	//
	//init static discover for example server
	//exampleStaticDiscover, e := discover.NewStaticDiscover("exampleproject", "examplegroup", "examplename", []string{"addr1","addr2"}, 9000, 10000, 8000)
	//if e != nil {
	//	return e
	//}
	//
	//init kubernetes discover for example server
	//exampleKubeDiscover, e := discover.NewKubernetesDiscover("exampleproject", "examplegroup", "examplename", "namespace", "fieldselector", "labelselector", 9000, 10000, 8000)
	//if e != nil {
	//	return e
	//}

	cgrpcc := config.GetCGrpcClientConfig().ClientConfig
	_ = cgrpcc //avoid unuse

	//init cgrpc client below
	//examplecgrpc, e = cgrpc.NewCGrpcClient(cgrpcc, examplediscover, model.Project, model.Group, model.Name, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	//         return e
	//}
	//ExampleCGrpcApi = example.NewExampleCGrpcClient(examplecgrpc)

	crpcc := config.GetCrpcClientConfig().ClientConfig
	_ = crpcc //avoid unuse

	//init crpc client below
	//examplecrpc, e = crpc.NewCrpcClient(crpcc, examplediscover, model.Project, model.Group, model.Name, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleCrpcApi = example.NewExampleCrpcClient(examplecrpc)

	webc := config.GetWebClientConfig().ClientConfig
	_ = webc //avoid unuse

	//init web client below
	//exampleweb, e = web.NewWebClient(webc, examplediscover, model.Project, model.Group, model.Name, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleWebApi = example.NewExampleWebClient(exampleweb)

	//DingTalk
	DingTalkStaticDiscover, e := discover.NewStaticDiscover("ali", "dingtalk", "oauth2", []string{"api.dingtalk.com"}, 0, 0, 0)
	if e != nil {
		return e
	}
	DingTalkWebClient, e = web.NewWebClient(webc, DingTalkStaticDiscover, model.Project, model.Group, model.Name, "ali", "dingtalk", "oauth2", &tls.Config{})
	if e != nil {
		return e
	}
	initDingTalk()

	//FeiShu
	FeiShuStaticDiscover, e := discover.NewStaticDiscover("bytedance", "feishu", "oauth2", []string{"open.feishu.cn"}, 0, 0, 0)
	if e != nil {
		return e
	}
	FeiShuWebClient, e = web.NewWebClient(webc, FeiShuStaticDiscover, model.Project, model.Group, model.Name, "bytedance", "feishu", "oauth2", &tls.Config{})
	if e != nil {
		return e
	}
	initFeiShu()

	return nil
}

func UpdateAppConfig(ac *config.AppConfig) {
	select {
	case trigerDingTalk <- nil:
	default:
	}
	select {
	case trigerFeiShu <- nil:
	default:
	}
}
