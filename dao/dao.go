package dao

import (
	"time"

	//"github.com/chenjie199234/admin/model"
	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/discover"
)

//var ExampleCGrpcApi example.ExampleCGrpcClient
//var ExampleCrpcApi example.ExampleCrpcClient
//var ExampleWebApi  example.ExampleWebClient

// NewApi create all dependent service's api we need in this program
func NewApi() error {
	var e error
	_ = e //avoid unuse

	//init discover for example server
	//examplediscover := NewDNSDiscover("examplegroup", "examplename", "examplename-headless.examplegroup", time.Second * 10, 9000, 10000, 8000)

	cgrpcc := GetCGrpcClientConfig()
	_ = cgrpcc //avoid unuse

	//init cgrpc client below
	//examplecgrpc, e = cgrpc.NewCGrpcClient(cgrpcc, examplediscover, model.Group, model.Name, "examplegroup", "examplename", nil)
	//if e != nil {
	//         return e
	//}
	//ExampleCGrpcApi = example.NewExampleCGrpcClient(examplecgrpc)

	crpcc := GetCrpcClientConfig()
	_ = crpcc //avoid unuse

	//init crpc client below
	//examplecrpc, e = crpc.NewCrpcClient(crpcc, examplediscover, model.Group, model.Name, "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleCrpcApi = example.NewExampleCrpcClient(examplecrpc)

	webc := GetWebClientConfig()
	_ = webc //avoid unuse

	//init web client below
	//exampleweb, e = web.NewWebClient(webc, examplediscover, model.Group, model.Name, "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleWebApi = example.NewExampleWebClient(exampleweb)

	return nil
}

func UpdateAPI(ac *config.AppConfig) {

}

func GetCGrpcClientConfig() *cgrpc.ClientConfig {
	gc := config.GetCGrpcClientConfig()
	return &cgrpc.ClientConfig{
		ConnectTimeout: time.Duration(gc.ConnectTimeout),
		GlobalTimeout:  time.Duration(gc.GlobalTimeout),
		HeartProbe:     time.Duration(gc.HeartProbe),
	}
}

func GetCrpcClientConfig() *crpc.ClientConfig {
	rc := config.GetCrpcClientConfig()
	return &crpc.ClientConfig{
		ConnectTimeout: time.Duration(rc.ConnectTimeout),
		GlobalTimeout:  time.Duration(rc.GlobalTimeout),
		HeartProbe:     time.Duration(rc.HeartProbe),
	}
}

func GetWebClientConfig() *web.ClientConfig {
	wc := config.GetWebClientConfig()
	return &web.ClientConfig{
		ConnectTimeout: time.Duration(wc.ConnectTimeout),
		GlobalTimeout:  time.Duration(wc.GlobalTimeout),
		IdleTimeout:    time.Duration(wc.IdleTimeout),
		HeartProbe:     time.Duration(wc.HeartProbe),
		MaxHeader:      2048,
	}
}
