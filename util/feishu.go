package util

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/cerror"
)

type getFeiShuUserTokenReq struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}
type getFeiShuUserTokenResp struct {
	Code int32                   `json:"code"`
	Msg  string                  `json:"msg"`
	Data *getFeiShuUserTokenData `json:"data"`
}
type getFeiShuUserTokenData struct {
	UserAccessToken string `json:"access_token"`
}
type getFeiShuUserInfoResp struct {
	Code int32                  `json:"code"`
	Msg  string                 `json:"msg"`
	Data *getFeiShuUserInfoData `json:"data"`
}
type getFeiShuUserInfoData struct {
	UserName string `json:"name"`
	UserID   string `json:"user_id"`
	Mobile   string `json:"mobile"`
}

func GetFeiShuOAuth2(ctx context.Context, code string) (username string, mobile string, e error) {
	//step1 get user token
	//https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/authen-v1/oidc-access_token/create?appId=cli_a596bbd826b8100d
	var usertoken string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json; charset=utf-8")
		header.Set("Authorization", "Bearer "+dao.FeiShuAppToken)
		req := &getFeiShuUserTokenReq{
			GrantType: "authorization_code",
			Code:      code,
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.FeiShuWebClient.Post(ctx, "/open-apis/authen/v1/oidc/access_token", "", header, nil, reqbody)
		if err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.usertoken] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.usertoken] read response body failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getFeiShuUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.usertoken] response body decode failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeCError(r.Code, 500, r.Msg)
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.usertoken] failed", slog.String("code", code), slog.String("error", e.Error()))
			return
		}
		usertoken = r.Data.UserAccessToken
	}
	//step2 get user info
	//https://open.feishu.cn/document/server-docs/authentication-management/login-state-management/get
	{
		header := make(http.Header)
		header.Set("Authorization", "Bearer "+usertoken)
		resp, err := dao.FeiShuWebClient.Get(ctx, "/open-apis/authen/v1/user_info", "", header, nil)
		if err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.userinfo] call failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.userinfo] read response body failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		r := &getFeiShuUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.userinfo] response body decode failed", slog.String("code", code), slog.String("error", err.Error()))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeCError(r.Code, 500, r.Msg)
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.userinfo] failed", slog.String("code", code), slog.String("error", e.Error()))
			return
		}
		username = r.Data.UserName
		if r.Data.Mobile == "" {
			e = ecode.ErrPermission
			slog.ErrorContext(ctx, "[GetFeiShuOAuth2.userinfo] missing mobile", slog.String("code", code), slog.String("user_name", username))
			return
		}
		mobile = r.Data.Mobile
	}
	return
}
