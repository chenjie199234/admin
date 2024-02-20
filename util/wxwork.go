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
	"github.com/chenjie199234/Corelib/util/common"
	"github.com/chenjie199234/Corelib/util/egroup"
)

type getWXWorkUserBaseInfoResp struct {
	Code       int32  `json:"errcode"`
	Msg        string `json:"errmsg"`
	UserID     string `json:"userid"`
	UserTicket string `json:"user_ticket"`
}
type getWXWorkUserMoreInfoResp struct {
	Code     int32  `json:"errcode"`
	Msg      string `json:"errmsg"`
	UserName string `json:"name"`
	Mobile   string `json:"mobile"`
}

func GetWXWorkOAuth2(ctx context.Context, code string) (username string, mobile string, e error) {
	//step1 get baseinfo
	//https://developer.work.weixin.qq.com/document/path/98176
	var userid string
	var userticket string
	{
		query := "access_token=" + dao.WXWorkAccessToken + "&code=" + code
		resp, err := dao.WXWorkWebClient.Get(ctx, "/cgi-bin/auth/getuserinfo", query, nil, nil)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &getWXWorkUserBaseInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		if r.Code != 0 {
			e = cerror.MakeError(r.Code, 500, r.Msg)
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] failed", log.String("code", code), log.CError(e))
			return
		}
		if r.UserID == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] doesn't delong to this corp", log.String("code", code), log.CError(e))
			return
		}
		if r.UserTicket == "" {
			e = ecode.ErrPermission
			log.Error(ctx, "[GetWXWorkOAuth2.baseinfo] can't get user ticket in wxwork", log.String("code", code), log.CError(e))
			return
		}
		userid = r.UserID
		userticket = r.UserTicket
	}
	//step2 get moreinfo
	eg := egroup.GetGroup(ctx)
	eg.Go(func(gctx context.Context) error {
		//https://developer.work.weixin.qq.com/document/path/90196
		query := "access_token=" + dao.WXWorkAccessToken + "&userid=" + userid
		resp, err := dao.WXWorkWebClient.Get(ctx, "/cgi-bin/user/get", query, nil, nil)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.username] call failed", log.String("code", code), log.CError(err))
			return err
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.username] read response body failed", log.String("code", code), log.CError(err))
			return err
		}
		r := &getWXWorkUserMoreInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.username] response body decode failed", log.String("code", code), log.CError(err))
			return err
		}
		if r.Code != 0 {
			err = cerror.MakeError(r.Code, 500, r.Msg)
			log.Error(ctx, "[GetWXWorkOAuth2.username] failed", log.String("code", code), log.CError(err))
			return err
		}
		if r.UserName == "" || r.UserName == userid {
			log.Error(ctx, "[GetWXWorkOAuth2.username] missing user name", log.String("code", code))
			return ecode.ErrPermission
		}
		username = r.UserName
		return nil
	})
	eg.Go(func(gctx context.Context) error {
		//https://developer.work.weixin.qq.com/document/path/95833
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		body := "{\"user_ticket\":\"" + userticket + "\"}"
		resp, err := dao.WXWorkWebClient.Post(ctx, "/cgi-bin/auth/getuserdetail", "access_token="+dao.WXWorkAccessToken, header, nil, common.STB(body))
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.usermobile] call failed", log.String("code", code), log.CError(err))
			return err
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.usermobile] read response body failed", log.String("code", code), log.CError(err))
			return err
		}
		r := &getWXWorkUserMoreInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetWXWorkOAuth2.usermobile] response body decode failed", log.String("code", code), log.CError(err))
			return err
		}
		if r.Code != 0 {
			err = cerror.MakeError(r.Code, 500, r.Msg)
			log.Error(ctx, "[GetWXWorkOAuth2.usermobile] failed", log.String("code", code), log.CError(err))
			return err
		}
		if r.Mobile == "" {
			log.Error(ctx, "[GetWXWorkOAuth2.usermobile] missing mobile", log.String("code", code))
			return ecode.ErrPermission
		}
		if r.Mobile[0] != '+' {
			mobile = "+86" + r.Mobile
		} else {
			mobile = r.Mobile
		}
		return nil
	})
	e = egroup.PutGroup(eg)
	return
}
