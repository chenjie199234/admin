package dao

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/cerror"
)

var FeiShuAppToken string
var trigerFeiShu chan *struct{}

// https://open.feishu.cn/document/server-docs/authentication-management/access-token/app_access_token_internal
func initFeiShu() {
	trigerFeiShu = make(chan *struct{}, 1)
	go func() {
		tmer := time.NewTimer(0)
		for {
			select {
			case <-tmer.C:
			case <-trigerFeiShu:
			}
			if tmer.Stop() {
				for len(tmer.C) > 0 {
					<-tmer.C
				}
			}
			if config.AC.Service.FeiShuOauth2 == "" || config.AC.Service.FeiShuAppID == "" || config.AC.Service.FeiShuAppSecret == "" {
				continue
			}
			r, e := getFeiShuAppToken()
			if e != nil {
				tmer.Reset(time.Millisecond * 500)
			} else {
				FeiShuAppToken = r.AppAccessToken
				tmer.Reset(time.Duration(r.ExpireIn-600) * time.Second)
			}
		}
	}()
}

type getFeiShuTokenReq struct {
	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`
}

type getFeiShuAppTokenResp struct {
	Code           int32  `json:"code"`
	Msg            string `json:"msg"`
	AppAccessToken string `json:"app_access_token"`
	ExpireIn       int64  `json:"expire"`
}

func getFeiShuAppToken() (*getFeiShuAppTokenResp, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json; charset=utf-8")
	req := &getFeiShuTokenReq{
		AppID:     config.AC.Service.FeiShuAppID,
		AppSecret: config.AC.Service.FeiShuAppSecret,
	}
	reqbody, _ := json.Marshal(req)
	resp, e := FeiShuWebClient.Post(context.Background(), "/open-apis/auth/v3/app_access_token/internal", "", header, nil, reqbody)
	if e != nil {
		slog.ErrorContext(nil, "[getFeiShuAppToken] call failed", slog.String("error", e.Error()))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		slog.ErrorContext(nil, "[getFeiShuAppToken] read respone body failed", slog.String("error", e.Error()))
		return nil, e
	}
	r := &getFeiShuAppTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		slog.ErrorContext(nil, "[getFeiShuAppToken] response body decode failed", slog.String("error", e.Error()))
		return nil, e
	}
	if r.Code != 0 {
		e = cerror.MakeCError(r.Code, 500, r.Msg)
		slog.ErrorContext(nil, "[GetFeiShuAppToken] failed", slog.String("error", e.Error()))
		return nil, e
	}
	return r, nil
}
func RefreshFeiShuToken() {
	select {
	case trigerFeiShu <- nil:
	default:
	}
}
