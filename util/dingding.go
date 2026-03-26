package util

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/chenjie199234/Corelib/cerror"
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
	ErrCode     int32  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
	CorpID      string `json:"corpId"`
}
type getDingDingUserInfoResp struct {
	ErrCode         int32  `json:"errcode"`
	ErrMsg          string `json:"errmsg"`
	UserName        string `json:"nick"`
	Mobile          string `json:"mobile"`
	MobileStateCode string `json:"stateCode"`
}

func GetDingDingOAuth2(ctx context.Context, clientid, clientsecret, code string) (username, mobile string, e error) {
	//step1 get user token
	//https://open.dingtalk.com/document/development/obtain-user-token
	var usertoken string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		req := &getDingDingUserTokenReq{
			ClientID:     clientid,
			ClientSecret: clientsecret,
			Code:         code,
			GrantType:    "authorization_code",
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingDingWebClient.Post(ctx, "/v1.0/oauth2/userAccessToken", "", header, nil, reqbody, nil)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] read response failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getDingDingUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] decode response failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		if r.ErrCode != 0 {
			e = cerror.MakeCError(r.ErrCode, 500, r.ErrMsg)
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.usertoken] oauth2 service provider return error",
				slog.String("code", code), slog.String("error", e.Error()))
			return "", "", e
		}
		usertoken = r.AccessToken
	}

	//step2 get user info
	//https://open.dingtalk.com/document/development/dingtalk-retrieve-user-information
	{
		header := make(http.Header)
		header.Set("x-acs-dingtalk-access-token", usertoken)
		resp, err := dao.DingDingWebClient.Get(ctx, "/v1.0/contact/users/me", "", header, nil, nil)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] read response failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getDingDingUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] decode response failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		if r.ErrCode != 0 {
			e = cerror.MakeCError(r.ErrCode, 500, r.ErrMsg)
			slog.ErrorContext(ctx, "[GetDingDingOAuth2.userinfo] oauth2 service provider return error",
				slog.String("code", code), slog.String("error", e.Error()))
			return "", "", e
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
