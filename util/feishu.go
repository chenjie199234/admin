package util

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
)

type getFeiShuUserTokenReq struct {
	GrantType string `json:"grant_type"`
	Code      string `json:"code"`
}
type getFeiShuUserTokenResp struct {
	Code int64                   `json:"code"`
	Msg  string                  `json:"msg"`
	Data *getFeiShuUserTokenData `json:"data"`
}
type getFeiShuUserTokenData struct {
	UserAccessToken string `json:"access_token"`
}
type getFeiShuUserInfoResp struct {
	Code int64                  `json:"code"`
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
			log.Error(ctx, "[GetFeiShuOAuth2.usertoken] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetFeiShuOAuth2.usertoken] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getFeiShuUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetFeiShuOAuth2.usertoken] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeError(int32(r.Code), 500, r.Msg)
			log.Error(ctx, "[GetFeiShuOAuth2.usertoken] failed", log.String("code", code), log.CError(e))
			return
		}
		usertoken = r.Data.UserAccessToken
	}
	//step2 get user info
	{
		header := make(http.Header)
		header.Set("Authorization", "Bearer "+usertoken)
		resp, err := dao.FeiShuWebClient.Get(ctx, "/open-apis/authen/v1/user_info", "", header, nil)
		if err != nil {
			log.Error(ctx, "[GetFeiShuOAuth2.userinfo] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetFeiShuOAuth2.userinfo] read response body failed", log.String("code", code), log.CError(e))
			e = err
			return
		}
		r := &getFeiShuUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetFeiShuOAuth2.userinfo] response body decode failed", log.String("code", code), log.CError(e))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeError(int32(r.Code), 500, r.Msg)
			log.Error(ctx, "[GetFeiShuOAuth2.userinfo] failed", log.String("code", code), log.CError(e))
			return
		}
		username = r.Data.UserName
		if r.Data.Mobile == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetFeiShuOAuth2.userinfo] missing mobile", log.String("code", code), log.String("user_name", username))
			return
		}
		mobile = r.Data.Mobile
	}
	return
}
