package util

import (
	"context"
	"encoding/json"
	"io"

	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
)

type getWXWorkUserIDResp struct {
	Code   int32  `json:"errcode"`
	Msg    string `json:"errmsg"`
	UserID string `json:"userid"`
}
type getWXWorkUserInfoResp struct {
	Code     int32  `json:"errcode"`
	Msg      string `json:"errmsg"`
	UserName string `json:"name"`
	Mobile   string `json:"mobile"`
}

func GetWXWorkOAuth2(ctx context.Context, code string) (username string, mobile string, e error) {
	//step1 get userid
	//https://developer.work.weixin.qq.com/document/path/98176
	var userid string
	{
		query := "access_token=" + dao.WXWorkAccessToken + "&code=" + code
		resp, err := dao.WXWorkWebClient.Get(ctx, "/cgi-bin/auth/getuserinfo", query, nil, nil)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userid] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userid] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getWXWorkUserIDResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userid] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeError(r.Code, 500, r.Msg)
			log.Error(ctx, "[GetWXWorkOAuth2.userid] failed", log.String("code", code), log.CError(e))
			return
		}
		if r.UserID == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetWXWorkOAuth2.userid] doesn't delong to this corp", log.String("code", code), log.CError(e))
			return
		}
		userid = r.UserID
	}
	//step2 get userinfo
	//https://developer.work.weixin.qq.com/document/path/90196
	{
		query := "access_token=" + dao.WXWorkAccessToken + "&userid=" + userid
		resp, err := dao.WXWorkWebClient.Get(ctx, "/cgi-bin/user/get", query, nil, nil)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getWXWorkUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeError(r.Code, 500, r.Msg)
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] failed", log.String("code", code), log.CError(e))
			return
		}
		if r.UserName == "" || r.UserName == userid {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] missing user name", log.String("code", code))
			return
		}
		if r.Mobile == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetWXWorkOAuth2.userinfo] missing mobile", log.String("code", code), log.String("user_name", r.UserName))
			return
		}
		username = r.UserName
		mobile = r.Mobile
	}
	return
}
