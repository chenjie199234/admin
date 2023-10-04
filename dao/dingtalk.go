package dao

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/web"
)

var DingTalkToken string
var trigerDingTalk chan *struct{}

func initDingTalk() {
	trigerDingTalk = make(chan *struct{}, 1)
	tmer := time.NewTimer(0)
	go func() {
		for {
			select {
			case <-tmer.C:
			case <-trigerDingTalk:
			}
			if tmer.Stop() {
				for len(tmer.C) > 0 {
					<-tmer.C
				}
			}
			if config.AC.Service.DingTalkOauth2 == "" || config.AC.Service.DingTalkAppKey == "" || config.AC.Service.DingTalkAppSecret == "" {
				DingTalkToken = ""
				continue
			}
			r, e := getDingTalkToken()
			if e != nil {
				tmer.Reset(time.Millisecond * 500)
			} else {
				DingTalkToken = r.AccessToken
				tmer.Reset(time.Duration(r.ExpireIn-600) * time.Second)
			}
		}
	}()
}

type getDingTalkTokenReq struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
type getDingTalkTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
}

func getDingTalkToken() (*getDingTalkTokenResp, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	req := &getDingTalkTokenReq{
		AppKey:    config.AC.Service.DingTalkAppKey,
		AppSecret: config.AC.Service.DingTalkAppSecret,
	}
	reqbody, _ := json.Marshal(req)
	resp, e := DingTalkWebClient.Post(web.WithForceAddr(context.Background(), "api.dingtalk.com"), "/v1.0/oauth2/accessToken", "", header, nil, reqbody)
	if e != nil {
		log.Error(nil, "[getDingTalkToken] call failed", log.CError(e))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		log.Error(nil, "[getDingTalkToken] read response body failed", log.CError(e))
		return nil, e
	}
	r := &getDingTalkTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		log.Error(nil, "[getDingTalkToken] response body decode failed", log.CError(e))
		return nil, e
	}
	return r, nil
}
