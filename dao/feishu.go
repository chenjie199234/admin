package dao

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"

	"github.com/chenjie199234/admin/config"
)

var FeiShuAppToken string
var trigerFeiShu chan *struct{}

func initFeiShu() {
	trigerFeiShu = make(chan *struct{})
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
	Code           int64  `json:"code"`
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
		log.Error(nil, "[getFeiShuAppToken] call failed", log.CError(e))
		return nil, e
	}
	defer resp.Body.Close()
	respbody, e := io.ReadAll(resp.Body)
	if e != nil {
		log.Error(nil, "[getFeiShuAppToken] read respone body failed", log.CError(e))
		return nil, e
	}
	r := &getFeiShuAppTokenResp{}
	if e = json.Unmarshal(respbody, r); e != nil {
		log.Error(nil, "[getFeiShuAppToken] response body decode failed", log.CError(e))
		return nil, e
	}
	if r.Code != 0 {
		e = cerror.MakeError(int32(r.Code), 500, r.Msg)
		log.Error(nil, "[GetFeiShuAppToken] failed", log.CError(e))
		return nil, e
	}
	return r, nil
}
