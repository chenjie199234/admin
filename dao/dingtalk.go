package dao

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chenjie199234/admin/config"

	"github.com/chenjie199234/Corelib/log"
)

var DingTalkCorpToken string
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
				DingTalkCorpToken = ""
				continue
			}
			r, e := getDingTalkCorpToken()
			if e != nil {
				tmer.Reset(time.Millisecond * 500)
			} else {
				DingTalkCorpToken = r.AccessToken
				tmer.Reset(time.Duration(r.ExpireIn-600) * time.Second)
			}
		}
	}()
}

type getDingTalkCorpTokenReq struct {
	AppKey    string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}
type getDingTalkCorpTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
}

func getDingTalkCorpToken() (*getDingTalkCorpTokenResp, error) {
	header := make(http.Header)
	header.Set("Content-Type", "application/json")
	req := &getDingTalkCorpTokenReq{
		AppKey:    config.AC.Service.DingTalkAppKey,
		AppSecret: config.AC.Service.DingTalkAppSecret,
	}
	reqbody, _ := json.Marshal(req)
	resp, e := DingTalkWebClient.Post(context.Background(), "/v1.0/oauth2/accessToken", "", header, nil, reqbody)
	if e != nil {
		log.Error(nil, "[getDingTalkCorpToken] call failed", log.CError(e))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		log.Error(nil, "[getDingTalkCorpToken] read response body failed", log.CError(e))
		return nil, e
	}
	r := &getDingTalkCorpTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		log.Error(nil, "[getDingTalkCorpToken] response body decode failed", log.CError(e))
		return nil, e
	}
	return r, nil
}
