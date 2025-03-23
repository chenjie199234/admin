package dao

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"time"

	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/cerror"
)

var WXWorkAccessToken string
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
		slog.ErrorContext(nil, "[getWXWorkAccessToken] call failed", slog.String("error", e.Error()))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		slog.ErrorContext(nil, "[getWXWorkAccessToken] read response body failed", slog.String("error", e.Error()))
		return nil, e
	}
	r := &getWXWorkAccessTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		slog.ErrorContext(nil, "[getWXWorkAccessToken] response body decode failed", slog.String("error", e.Error()))
		return nil, e
	}
	if r.Code != 0 {
		e = cerror.MakeCError(r.Code, 500, r.Msg)
		slog.ErrorContext(nil, "[getWXWorkAccessToken] failed", slog.String("error", e.Error()))
		return nil, e
	}
	return r, nil
}
func RefreshWXWorkToken() {
	select {
	case trigerWXWork <- nil:
	default:
	}
}
