package dao

import (
	"net"
	"time"

	//"github.com/chenjie199234/admin/api"
	//example "github.com/chenjie199234/admin/api/deps/example"
	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/cgrpc"
	"github.com/chenjie199234/Corelib/crpc"
	"github.com/chenjie199234/Corelib/web"
)

//var ExampleCGrpcApi example.ExampleCGrpcClient
//var ExampleCrpcApi example.ExampleCrpcClient
//var ExampleWebApi  example.ExampleWebClient

//NewApi create all dependent service's api we need in this program
func NewApi() error {
	var e error
	_ = e //avoid unuse

	cgrpcc := getCGrpcClientConfig()
	_ = cgrpcc //avoid unuse

	//init cgrpc client below
	//examplecgrpc e = cgrpc.NewCGrpcClient(cgrpcc, api.Group, api.Name, "examplegroup", "examplename")
	//if e != nil {
	//         return e
	//}
	//ExampleCGrpcApi = example.NewExampleCGrpcClient(examplecgrpc)

	crpcc := getCrpcClientConfig()
	_ = crpcc //avoid unuse

	//init crpc client below
	//examplecrpc, e = crpc.NewCrpcClient(crpcc, api.Group, api.Name, "examplegroup", "examplename")
	//if e != nil {
	// 	return e
	//}
	//ExampleCrpcApi = example.NewExampleCrpcClient(examplecrpc)

	webc := getWebClientConfig()
	_ = webc //avoid unuse

	//init web client below
	//exampleweb, e = web.NewWebClient(webc, api.Group, api.Name, "examplegroup", "examplename", "http://examplehost:exampleport")
	//if e != nil {
	// 	return e
	//}
	//ExampleWebApi = example.NewExampleWebClient(exampleweb)

	return nil
}

func getCGrpcClientConfig() *cgrpc.ClientConfig {
	gc := config.GetCGrpcClientConfig()
	return &cgrpc.ClientConfig{
		ConnectTimeout:   time.Duration(gc.ConnectTimeout),
		GlobalTimeout:    time.Duration(gc.GlobalTimeout),
		HeartPorbe:       time.Duration(gc.HeartProbe),
		Discover:         cgrpcDNS,
		DiscoverInterval: time.Second * 10,
	}
}

func cgrpcDNS(group, name string) (map[string]*cgrpc.RegisterData, error) {
	result := make(map[string]*cgrpc.RegisterData)
	addrs, e := net.LookupHost(name + "-service-headless." + group)
	if e != nil {
		return nil, e
	}
	for i := range addrs {
		addrs[i] = addrs[i] + ":10000"
	}
	dserver := make(map[string]*struct{})
	dserver["dns"] = nil
	for _, addr := range addrs {
		result[addr] = &cgrpc.RegisterData{DServers: dserver}
	}
	return result,nil
}

func getCrpcClientConfig() *crpc.ClientConfig {
	rc := config.GetCrpcClientConfig()
	return &crpc.ClientConfig{
		ConnectTimeout:   time.Duration(rc.ConnectTimeout),
		GlobalTimeout:    time.Duration(rc.GlobalTimeout),
		HeartPorbe:       time.Duration(rc.HeartProbe),
		Discover:         crpcDNS,
		DiscoverInterval: time.Second * 10,
	}
}

func crpcDNS(group, name string) (map[string]*crpc.RegisterData, error) {
	result := make(map[string]*crpc.RegisterData)
	addrs, e := net.LookupHost(name + "-service-headless." + group)
	if e != nil {
		return nil, e
	}
	for i := range addrs {
		addrs[i] = addrs[i] + ":9000"
	}
	dserver := make(map[string]*struct{})
	dserver["dns"] = nil
	for _, addr := range addrs {
		result[addr] = &crpc.RegisterData{DServers: dserver}
	}
	return result, nil
}

func getWebClientConfig() *web.ClientConfig {
	wc := config.GetWebClientConfig()
	return &web.ClientConfig{
		ConnectTimeout: time.Duration(wc.ConnectTimeout),
		GlobalTimeout:  time.Duration(wc.GlobalTimeout),
		IdleTimeout:    time.Duration(wc.IdleTimeout),
		HeartProbe:     time.Duration(wc.HeartProbe),
		MaxHeader:      1024,
	}
}