package util

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"
)

type getDingDingUserTokenReq struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Code         string `json:"code"`
	GrantType    string `json:"grantType"`
}
type getDingDingUserTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
	CorpID      string `json:"corpId"`
}
type getDingDingUserInfoResp struct {
	UnionID         string `json:"unionId"`
	UserName        string `json:"nick"`
	Mobile          string `json:"mobile"`
	MobileStateCode string `json:"stateCode"`
}

func GetDingDingOAuth2(ctx context.Context, code string) (username, mobile string, e error) {
	//step1 get user token
	//https://open.dingtalk.com/document/orgapp/obtain-user-token
	var usertoken string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		req := &getDingDingUserTokenReq{
			ClientID:     config.AC.Service.DingDingClientID,
			ClientSecret: config.AC.Service.DingDingClientSecret,
			Code:         code,
			GrantType:    "authorization_code",
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingDingWebClient.Post(ctx, "/v1.0/oauth2/userAccessToken", "", header, nil, reqbody)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] read response body failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getDingDingUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] response body decode failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		usertoken = r.AccessToken
	}

	//step2 get user info
	//https://open.dingtalk.com/document/orgapp/dingtalk-retrieve-user-information
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("x-acs-dingtalk-access-token", usertoken)
		resp, err := dao.DingDingWebClient.Get(ctx, "/v1.0/contact/users/me", "", header, nil)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] read response body failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getDingDingUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] response body deocde failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		username = r.UserName
		if r.MobileStateCode == "" || r.Mobile == "" {
			e = ecode.ErrPermission
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] missing mobile", slog.String("code", code), slog.String("user_name", username))
			return
		}
		mobile = "+" + r.MobileStateCode + r.Mobile
	}
	return
}
