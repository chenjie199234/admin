package util

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/chenjie199234/admin/config"
	"github.com/chenjie199234/admin/dao"
	"github.com/chenjie199234/admin/ecode"

	"github.com/chenjie199234/Corelib/cerror"
	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/egroup"
	"github.com/chenjie199234/Corelib/web"
)

type GetDingTalkUserTokenReq struct {
	AppKey    string `json:"clientId"`
	AppSecret string `json:"clientSecret"`
	Code      string `json:"code"`
	GrantType string `json:"grantType"`
}
type GetDingTalkUserTokenResp struct {
	AccessToken string `json:"accessToken"`
	ExpireIn    int64  `json:"expireIn"`
	CorpID      string `json:"corpId"`
}
type GetDingTalkUnionIDByUserTokenResp struct {
	UnionID string `json:"unionId"`
}
type GetDingTalkUserIDByUnionIDReq struct {
	UnionID string `json:"unionid"`
}
type GetDingTalkUserIDByUnionIDResp struct {
	ErrCode int64                           `json:"errcode"`
	ErrMsg  string                          `json:"errmsg"`
	Result  *GetDingTalkUserIDByUnionIDData `json:"result"`
}
type GetDingTalkUserIDByUnionIDData struct {
	ContactType int    `json:"contact_type"` //0 inner staff,1 outside contacter
	UserID      string `json:"userid"`
}

func GetDingTalkOAuth2(ctx context.Context, code string) (userid string, e error) {
	//step1 get user token
	var usertoken string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		req := &GetDingTalkUserTokenReq{
			AppKey:    config.AC.Service.DingTalkAppKey,
			AppSecret: config.AC.Service.DingTalkAppSecret,
			Code:      code,
			GrantType: "authorization_code",
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingTalkWebClient.Post(web.WithForceAddr(ctx, "api.dingtalk.com"), "/v1.0/oauth2/userAccessToken", "", header, nil, reqbody)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &GetDingTalkUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.usertoken] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		usertoken = r.AccessToken
	}

	//step2 get unionid
	var unionid string
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Set("x-acs-dingtalk-access-token", usertoken)
		resp, err := dao.DingTalkWebClient.Get(web.WithForceAddr(ctx, "api.dingtalk.com"), "/v1.0/contact/users/me", "", header, nil)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.unionid] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.unionid] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &GetDingTalkUnionIDByUserTokenResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.unionid] response body deocde failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		unionid = r.UnionID
	}

	//step3 get userid
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		header.Del("x-acs-dingtalk-access-token")
		req := &GetDingTalkUserIDByUnionIDReq{
			UnionID: unionid,
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingTalkWebClient.Post(web.WithForceAddr(ctx, "oapi.dingtalk.com"), "/topapi/user/getbyunionid", "access_token="+dao.DingTalkCorpToken, header, nil, reqbody)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userid] call failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userid] read response body failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		r := &GetDingTalkUserIDByUnionIDResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkOAuth2.userid] response body decode failed", log.String("code", code), log.CError(err))
			e = err
			return
		}
		//https://open.dingtalk.com/document/orgapp/query-a-user-by-the-union-id
		if r.ErrCode != 0 {
			if r.ErrCode == 60121 {
				e = ecode.ErrUserNotExist
			} else if r.ErrCode == -1 {
				e = ecode.ErrBusy
			} else {
				e = cerror.MakeError(int32(r.ErrCode), 500, r.ErrMsg)
			}
			log.Error(ctx, "[GetDingTalkOAuth2.userid] failed", log.String("code", code), log.CError(e))
			return
		}
		if r.Result.ContactType == 0 {
			userid = r.Result.UserID
		}
	}
	return
}

type GetDingTalkUserInfoReq struct {
	UserID   string `json:"userid"`
	Language string `json:"language"`
}
type GetDingTalkUserInfoResp struct {
	ErrCode int64                    `json:"errcode"`
	ErrMsg  string                   `json:"errmsg"`
	Result  *GetDingTalkUserInfoData `json:"result"`
}
type GetDingTalkUserInfoData struct {
	Name            string  `json:"name"`
	MobileStateCode string  `json:"state_code"`
	Mobile          string  `json:"mobile"`
	DepartmentIDs   []int64 `json:"dept_id_list"`
}
type GetDingTalkDepartmentReq struct {
	DepartmentID int64  `json:"dept_id"`
	Language     string `json:"language"`
}
type GetDingTalkDepartmentResp struct {
	ErrCode int64                      `json:"errcode"`
	ErrMsg  string                     `json:"errmsg"`
	Result  *GetDingTalkDepartmentData `json:"result"`
}
type GetDingTalkDepartmentData struct {
	Name string `json:"name"`
}

func GetDingTalkUserInfo(ctx context.Context, userid string) (name, mobile, department string, e error) {
	var departmentids []int64
	{
		header := make(http.Header)
		header.Set("Content-Type", "application/json")
		req := GetDingTalkUserInfoReq{
			UserID:   userid,
			Language: "zh_CN",
		}
		reqbody, _ := json.Marshal(req)
		resp, err := dao.DingTalkWebClient.Post(web.WithForceAddr(ctx, "oapi.dingtalk.com"), "/topapi/v2/user/get", "access_token="+dao.DingTalkCorpToken, header, nil, reqbody)
		if err != nil {
			log.Error(ctx, "[GetDingTalkUserInfo.user] call failed", log.String("oauth2_user_id", userid), log.CError(err))
			e = err
			return
		}
		defer resp.Body.Close()
		respbody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Error(ctx, "[GetDingTalkUserInfo.user] read response body failed", log.String("oauth2_user_id", userid), log.CError(err))
			e = err
			return
		}
		r := &GetDingTalkUserInfoResp{}
		if err = json.Unmarshal(respbody, r); err != nil {
			log.Error(ctx, "[GetDingTalkUserInfo.user] response body decode failed", log.String("oauth2_user_id", userid), log.CError(err))
			e = err
			return
		}
		//https://open.dingtalk.com/document/orgapp/query-user-details
		if r.ErrCode != 0 {
			if r.ErrCode == -1 {
				e = ecode.ErrBusy
			} else if r.ErrCode == 33012 {
				e = ecode.ErrUserNotExist
			} else {
				e = cerror.MakeError(int32(r.ErrCode), 500, r.ErrMsg)
			}
			log.Error(ctx, "[GetDingTalkUserInfo.user] failed", log.String("oauth2_user_id", userid), log.CError(e))
			return
		}
		name = r.Result.Name
		if r.Result.MobileStateCode != "" {
			mobile = "+" + r.Result.MobileStateCode + "-" + r.Result.Mobile
		} else {
			mobile = r.Result.Mobile
		}
		departmentids = r.Result.DepartmentIDs
	}
	departmentnames := make(map[int64]string, len(departmentids))
	eg := egroup.GetGroup(ctx)
	for _, v := range departmentids {
		departmentid := v
		departmentnames[departmentid] = ""
		eg.Go(func(gctx context.Context) error {
			header := make(http.Header)
			header.Set("Content-Type", "application/json")
			req := &GetDingTalkDepartmentReq{
				DepartmentID: departmentid,
				Language:     "zh_CN",
			}
			reqbody, _ := json.Marshal(req)
			resp, err := dao.DingTalkWebClient.Post(web.WithForceAddr(ctx, "oapi.dingtalk.com"), "/topapi/v2/department/get", "access_token="+dao.DingTalkCorpToken, header, nil, reqbody)
			if err != nil {
				log.Error(ctx, "[GetDingTalkUserInfo.department] call failed",
					log.String("oauth2_user_id", userid),
					log.Int64("oauth2_department_id", departmentid),
					log.CError(err))
				return err
			}
			defer resp.Body.Close()
			respbody, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Error(ctx, "[GetDingTalkUserInfo.department] read response body failed",
					log.String("oauth2_user_id", userid),
					log.Int64("oauth2_department_id", departmentid),
					log.CError(err))
				return err
			}
			r := &GetDingTalkDepartmentResp{}
			if err = json.Unmarshal(respbody, r); err != nil {
				log.Error(ctx, "[GetDingTalkUserInfo.department] response body decode failed",
					log.String("oauth2_user_id", userid),
					log.Int64("oauth2_department_id", departmentid),
					log.CError(err))
				return err
			}
			//https://open.dingtalk.com/document/orgapp/query-department-details0-v2
			if r.ErrCode != 0 {
				if r.ErrCode == -1 {
					err = ecode.ErrBusy
				} else {
					err = cerror.MakeError(int32(r.ErrCode), 500, r.ErrMsg)
				}
				log.Error(ctx, "[GetDingTalkUserInfo.department] failed",
					log.String("oauth2_user_id", userid),
					log.Int64("oauth2_department_id", departmentid),
					log.CError(err))
				return err
			}
			departmentnames[departmentid] = r.Result.Name
			return nil
		})
	}
	e = egroup.PutGroup(eg)
	if e != nil {
		return
	}
	for i, departmentid := range departmentids {
		if i != 0 {
			department += "-"
		}
		department += departmentnames[departmentid]
	}
	return
}
