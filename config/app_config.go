package config

import (
	"log/slog"
	"os"

	publicmids "github.com/chenjie199234/Corelib/mids"
	"github.com/chenjie199234/Corelib/util/ctime"
)

// AppConfig can hot update
// this is the config used for this app
type AppConfig struct {
	HandlerTimeout map[string]map[string]ctime.Duration `json:"handler_timeout"`  //first key path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite map[string]map[string]string         `json:"web_path_rewrite"` //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate    publicmids.MultiPathRateConfigs      `json:"handler_rate"`     //key:path
	Accesses       publicmids.MultiPathAccessConfigs    `json:"accesses"`         //key:path
	TokenSecret    string                               `json:"token_secret"`     //if don't need token check,this can be ingored
	Service        *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
	TokenExpire ctime.Duration `json:"token_expire"`

	//https://login.dingtalk.com/oauth2/auth?redirect_uri={REDIRECT_URI}&response_type=code&client_id={CLIENTID}&scope=openid&state=DingDing&prompt=consent
	DingDingOauth2       string `json:"dingding_oauth2"`
	DingDingClientID     string `json:"dingding_client_id"`
	DingDingClientSecret string `json:"dingding_client_secret"`

	//https://accounts.feishu.cn/open-apis/authen/v1/authorize?redirect_uri={REDIRECT_URI}&response_type=code&client_id={CLIENT_ID}&state=FeiShu&scope=contact:user.phone:readonly
	FeiShuOauth2       string `json:"feishu_oauth2"`
	FeiShuClientID     string `json:"feishu_client_id"`
	FeiShuClientSecret string `json:"feishu_client_secret"`

	//https://open.weixin.qq.com/connect/oauth2/authorize?redirect_uri={REDIRECT_URI}&appid={CORPID}&response_type=code&scope=snsapi_privateinfo&state=WXWork&agentid={AGENTID}#wechat_redirect
	WXWorkOauth2     string `json:"wxwork_oauth2"`
	WXWorkCorpID     string `json:"wxwork_corp_id"`
	WXWorkCorpSecret string `json:"wxwork_corp_secret"`
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
	if ac.Service.DingDingOauth2 == "" && ac.Service.FeiShuOauth2 == "" && ac.Service.WXWorkOauth2 == "" {
		slog.Warn("[config.validateAppConfig] no oauth2 service,only root can login by password")
	}
	if ac.Service.DingDingOauth2 != "" && (ac.Service.DingDingClientID == "" || ac.Service.DingDingClientSecret == "") {
		slog.Error("[config.validateAppConfig] missing dingding_client_id or dingding_client_secret")
		os.Exit(1)
	}
	if ac.Service.FeiShuOauth2 != "" && (ac.Service.FeiShuClientID == "" || ac.Service.FeiShuClientSecret == "") {
		slog.Error("[config.validateAppConfig] missing feishu_client_id or feishu_client_secret")
		os.Exit(1)
	}
	if ac.Service.WXWorkOauth2 != "" && (ac.Service.WXWorkCorpID == "" || ac.Service.WXWorkCorpSecret == "") {
		slog.Error("[config.validateAppConfig] missing wxwork_corp_id or wxwork_corp_secret")
		os.Exit(1)
	}
}

// AC -
var AC *AppConfig
