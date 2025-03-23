package dao

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/chenjie199234/admin/config"
)

var DingDingToken string
var trigerDingDing chan *struct{}

// https://open.dingtalk.com/document/orgapp/obtain-the-access_token-of-an-internal-app?spm=ding_open_doc.document.0.0.10686dc3KUu0Es
func initDingDing() {
	trigerDingDing = make(chan *struct{}, 1)
	go func() {
		tmer := time.NewTimer(0)
		for {
			select {
			case <-tmer.C:
			case <-trigerDingDing:
			}
			tmer.Stop()
			r, e := getDingDingToken()
			if e != nil {
				tmer.Reset(time.Millisecond * 500)
			} else if r == nil {
				DingDingToken = ""
				continue
			} else {
				DingDingToken = r.AccessToken
				tmer.Reset(time.Duration(r.ExpireIn-600) * time.Second)
			}
		}
	}()
}

type getDingDingTokenReq struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
type getDingDingTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
}

func getDingDingToken() (*getDingDingTokenResp, error) {
	c := config.AC.Service
	if c.DingDingOauth2 == "" || c.DingDingClientID == "" || c.DingDingClientSecret == "" {
		return nil, nil
	}
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	req := &getDingDingTokenReq{
		AppKey:    c.DingDingClientID,
		AppSecret: c.DingDingClientSecret,
	}
	reqbody, _ := json.Marshal(req)
	resp, e := DingDingWebClient.Post(context.Background(), "/v1.0/oauth2/accessToken", "", header, nil, reqbody)
	if e != nil {
		slog.ErrorContext(nil, "[getDingDingToken] call failed", slog.String("error", e.Error()))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		slog.ErrorContext(nil, "[getDingDingToken] read response body failed", slog.String("error", e.Error()))
		return nil, e
	}
	r := &getDingDingTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		slog.ErrorContext(nil, "[getDingDingToken] response body decode failed", slog.String("error", e.Error()))
		return nil, e
	}
	return r, nil
}
func RefreshDingDingToken() {
	select {
	case trigerDingDing <- nil:
	default:
	}
}
