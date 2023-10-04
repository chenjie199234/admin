package util

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/log"
)

type getDingTalkUserTokenReq struct {
	AppKey    string `json:"clientId"`
	AppSecret string `json:"clientSecret"`
	Code      string `json:"code"`
	GrantType string `json:"grantType"`
}
type getDingTalkUserTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
	CorpID      string `json:"corpId"`
}
type getDingTalkUnionIDByUserTokenResp struct {
	UnionID         string `json:"unionId"`
	UserName        string `json:"nick"`
	Mobile          string `json:"mobile"`
	MobileStateCode string `json:"stateCode"`
}
type getDingTalkUserIDByUnionIDReq struct {
	UnionID string `json:"unionid"`
}
type getDingTalkUserIDByUnionIDResp struct {
	ErrCode int64                           `json:"errcode"`
	ErrMsg  string                          `json:"errmsg"`
	Result  *getDingTalkUserIDByUnionIDData `json:"result"`
}
type getDingTalkUserIDByUnionIDData struct {
	ContactType int    `json:"contact_type"` //0 inner staff,1 outside contacter
	UserID      string `json:"userid"`
}

func GetDingTalkOAuth2(ctx context.Context, code string) (username, mobile string, e error) {
	//step1 get user token
	var usertoken string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		req := &getDingTalkUserTokenReq{
			AppKey:    config.AC.Service.DingTalkAppKey,
			AppSecret: config.AC.Service.DingTalkAppSecret,
			Code:      code,
			GrantType: "authorization_code",
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingTalkWebClient.Post(ctx, "/v1.0/oauth2/userAccessToken", "", header, nil, reqbody)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getDingTalkUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		usertoken = r.AccessToken
	}

	//step2 get user info
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("x-acs-dingtalk-access-token", usertoken)
		resp, err := dao.DingTalkWebClient.Get(ctx, "/v1.0/contact/users/me", "", header, nil)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userinfo] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userinfo] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getDingTalkUnionIDByUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userinfo] response body deocde failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		username = r.UserName
		if r.MobileStateCode == "" || r.Mobile == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetDingTalkOAuth2.userinfo] missing mobile", log.String("code", code), log.String("user_name", username))
			return
		}
		mobile = "+" + r.MobileStateCode + r.Mobile
	}
	return
}
