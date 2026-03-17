package dao

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/discover"
	"github.com/chenjie199234/Corelib/web"
	// "github.com/chenjie199234/Corelib/cgrpc"
	// "github.com/chenjie199234/Corelib/crpc"
)

// var ExampleCGrpcApi example.ExampleCGrpcClient
// var ExampleCrpcApi example.ExampleCrpcClient
// var ExampleWebApi  example.ExampleWebClient
var DingDingWebClient *web.WebClient
var FeiShuWebClient *web.WebClient
var WXWorkWebClient *web.WebClient

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
	//examplecgrpc, e = cgrpc.NewCGrpcClient(cgrpcc, examplediscover, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	//         return e
	//}
	//ExampleCGrpcApi = example.NewExampleCGrpcClient(examplecgrpc)

	crpcc := config.GetCrpcClientConfig().ClientConfig
	_ = crpcc //avoid unuse

	//init crpc client below
	//examplecrpc, e = crpc.NewCrpcClient(crpcc, examplediscover, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleCrpcApi = example.NewExampleCrpcClient(examplecrpc)

	webc := config.GetWebClientConfig().ClientConfig
	_ = webc //avoid unuse

	//init web client below
	//exampleweb, e = web.NewWebClient(webc, examplediscover, "exampleproject", "examplegroup", "examplename", nil)
	//if e != nil {
	// 	return e
	//}
	//ExampleWebApi = example.NewExampleWebClient(exampleweb)

	//DingTalk
	DingTalkStaticDiscover, e := discover.NewStaticDiscover("ali", "dingtalk", "oauth2", []string{"api.dingtalk.com"}, 0, 0, 0)
	if e != nil {
		return e
	}
	DingDingWebClient, e = web.NewWebClient(webc, DingTalkStaticDiscover, "ali", "dingtalk", "oauth2", &tls.Config{})
	if e != nil {
		return e
	}

	//FeiShu
	FeiShuStaticDiscover, e := discover.NewStaticDiscover("bytedance", "feishu", "oauth2", []string{"open.feishu.cn"}, 0, 0, 0)
	if e != nil {
		return e
	}
	FeiShuWebClient, e = web.NewWebClient(webc, FeiShuStaticDiscover, "bytedance", "feishu", "oauth2", &tls.Config{})
	if e != nil {
		return e
	}

	//WXWork
	WXWorkStaticDiscover, e := discover.NewStaticDiscover("tencent", "wxwork", "oauth2", []string{"qyapi.weixin.qq.com"}, 0, 0, 0)
	if e != nil {
		return e
	}
	WXWorkWebClient, e = web.NewWebClient(webc, WXWorkStaticDiscover, "tencent", "wxwork", "oauth2", &tls.Config{})
	if e != nil {
		return e
	}
	initWXWork()

	return nil
}

func UpdateAppConfig(ac *config.AppConfig) {
	RefreshWXWorkToken()
}

var WXWorkAccessToken string
var curIdSecret string
var trigerWXWork chan *struct{}

// https://developer.work.weixin.qq.com/document/path/91039
func initWXWork() {
	trigerWXWork = make(chan *struct{}, 1)
	go func() {
		tmer := time.NewTimer(0)
		for {
			select {
			case <-tmer.C:
			case <-trigerWXWork:
			}
			tmer.Stop()
			r, e := getWXWorkAccessToken()
			if e != nil {
				tmer.Reset(time.Millisecond * 500)
			} else if r == nil {
				WXWorkAccessToken = ""
				continue
			} else {
				WXWorkAccessToken = r.AccessToken
				tmer.Reset(time.Duration(r.ExpireIn) * time.Second)
			}
			select {
			case <-trigerWXWork:
			default:
			}
		}
	}()
}

type getWXWorkAccessTokenResp struct {
	Code        int32  `json:"errcode"`
	Msg         string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpireIn    int64  `json:"expires_in"`
}

func getWXWorkAccessToken() (*getWXWorkAccessTokenResp, error) {
	c := config.AC.Service
	if c.WXWorkOauth2 == "" || c.WXWorkCorpID == "" || c.WXWorkCorpSecret == "" {
		return nil, nil
	}
	query := "corpid=" + c.WXWorkCorpID + "&corpsecret=" + c.WXWorkCorpSecret
	resp, e := WXWorkWebClient.Get(context.Background(), "/cgi-bin/gettoken", query, nil, nil)
	if e != nil {
		slog.Error("[getWXWorkAccessToken] call failed", slog.String("error", e.Error()))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		slog.Error("[getWXWorkAccessToken] read response body failed", slog.String("error", e.Error()))
		return nil, e
	}
	r := &getWXWorkAccessTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		slog.Error("[getWXWorkAccessToken] response body decode failed", slog.String("error", e.Error()))
		return nil, e
	}
	if r.Code != 0 {
		e = cerror.MakeCError(r.Code, 500, r.Msg)
		slog.Error("[getWXWorkAccessToken] failed", slog.String("error", e.Error()))
		return nil, e
	}
	curIdSecret = c.WXWorkCorpID + c.WXWorkCorpSecret
	return r, nil
}
func RefreshWXWorkToken() {
	c := config.AC.Service
	if curIdSecret == c.WXWorkCorpID+c.WXWorkCorpSecret {
		return
	}
	select {
	case trigerWXWork <- nil:
	default:
	}
}
