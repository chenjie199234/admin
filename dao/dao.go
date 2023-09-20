package dao

import (
	"github.com/chenjie199234/admin/config"
	// "github.com/chenjie199234/admin/model"
	// "github.com/chenjie199234/Corelib/discover"
	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
	// "github.com/chenjie199234/Corelib/web"
)

//var ExampleCGrpcApi example.ExampleCGrpcClient
//var ExampleCrpcApi example.ExampleCrpcClient
//var ExampleWebApi  example.ExampleWebClient

// NewApi create all dependent service's api we need in this program
func NewApi() error {
	//init discover for example server
	//examplediscover, e := discover.NewDNSDiscover("exampleproject", "examplegroup", "examplename", "exampleproject-examplegroup.examplename-headless", time.Second * 10, 9000, 10000, 8000)
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

	return nil
}

func UpdateAPI(ac *config.AppConfig) {

}
