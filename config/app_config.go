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
	HandlerTimeout map[string]map[string]ctime.Duration `json:"handler_timeout"`  //first key handler path,second key method(GET,POST,PUT,PATCH,DELETE,CRPC,GRPC),value timeout
	WebPathRewrite map[string]map[string]string         `json:"web_path_rewrite"` //first key method(GET,POST,PUT,PATCH,DELETE),second key origin url,value new url
	HandlerRate    publicmids.MultiPathRateConfigs      `json:"handler_rate"`     //key:path
	Accesses       publicmids.MultiPathAccessConfigs    `json:"accesses"`         //key:path
	TokenSecret    string                               `json:"token_secret"`
	Service        *ServiceConfig                       `json:"service"`
}
type ServiceConfig struct {
	//add your config here
	TokenExpire ctime.Duration `json:"token_expire"`

	//https://login.dingtalk.com/oauth2/auth?redirect_uri={REDIRECT_URI}&response_type=code&client_id={CLIENTID}&scope=openid&state=DingDing&prompt=consent
	DingDingOauth2       string `json:"dingding_oauth2"`
	DingDingClientID     string `json:"dingding_client_id"`
	DingDingClientSecret string `json:"dingding_client_secret"`

	//https://open.feishu.cn/open-apis/authen/v1/authorize?redirect_uri={REDIRECT_URI}&app_id={APPID}&state=FeiShu&scope=contact:user.employee_id:readonly%20contact:user.phone:readonly
	FeiShuOauth2    string `json:"feishu_oauth2"`
	FeiShuAppID     string `json:"feishu_app_id"`
	FeiShuAppSecret string `json:"feishu_app_secret"`

	//https://open.weixin.qq.com/connect/oauth2/authorize?redirect_uri={REDIRECT_URI}&appid={CORPID}&response_type=code&scope=snsapi_privateinfo&state=WXWork&agentid={AGENTID}#wechat_redirect
	WXWorkOauth2     string `json:"wxwork_oauth2"`
	WXWorkCorpID     string `json:"wxwork_corp_id"`
	WXWorkCorpSecret string `json:"wxwork_corp_secret"`
}

// every time update AppConfig will call this function
func validateAppConfig(ac *AppConfig) {
	if ac.Service.DingDingOauth2 == "" && ac.Service.FeiShuOauth2 == "" && ac.Service.WXWorkOauth2 == "" {
		slog.WarnContext(nil, "[config.validateAppConfig] no oauth2 service,only root can login by password")
	}
	if ac.Service.DingDingOauth2 != "" && (ac.Service.DingDingClientID == "" || ac.Service.DingDingClientSecret == "") {
		slog.ErrorContext(nil, "[config.validateAppConfig] missing dingding_client_id or dingding_client_secret")
		os.Exit(1)
	}
	if ac.Service.FeiShuOauth2 != "" && (ac.Service.FeiShuAppID == "" || ac.Service.FeiShuAppSecret == "") {
		slog.ErrorContext(nil, "[config.validateAppConfig] missing feishu_app_id or feishu_app_secret")
		os.Exit(1)
	}
	if ac.Service.WXWorkOauth2 != "" && (ac.Service.WXWorkCorpID == "" || ac.Service.WXWorkCorpSecret == "") {
		slog.ErrorContext(nil, "[config.validateAppConfig] missing wxwork_corp_id or wxwork_corp_secret")
		os.Exit(1)
	}
}

// AC -
var AC *AppConfig
