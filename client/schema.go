package client

import (
	"fmt"
	"time"
)

//ResAccessToken struct
type accessToken struct {
	ExpireMessage
	AccessToken string `json:"access_token"`
}

// resTicket 请求jsapi_tikcet返回结果
type resTicket struct {
	ExpireMessage
	Ticket string `json:"ticket"`
}

type CommonResp struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (msg *CommonResp) Error() error {
	if msg.ErrCode != 0 {
		return fmt.Errorf("result code=%d: %s", msg.ErrCode, msg.ErrMsg)
	}
	return nil
}

type ExpireMessage struct {
	CommonResp
	ExpiresIn int64 `json:"expires_in"`
	LastUpTs  int64 `json:"update_in"`
}

func (msg *ExpireMessage) valid() bool {
	return msg.ErrCode == 0 && time.Now().Unix() < msg.LastUpTs+msg.ExpiresIn-60
}

// 素材返回结果
type Material struct {
	CommonResp
	Type      MediaType `json:"type"`
	MediaID   string    `json:"media_id"`
	CreatedAt int64     `json:"created_at"`
	URL       string    `json:"url"`
}

//UserInfo 用户授权获取到用户信息
type UserInfo struct {
	CommonResp

	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

type QrCode struct {
	CommonResp
	ExpiresSeconds int64  `json:"expire_seconds"`
	Ticket         string `json:"ticket"`
	Url            string `json:"url"`
}

// UserAccessToken 获取用户授权access_token的返回结果
type UserAccessToken struct {
	ExpireMessage

	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

// UserSession 登录凭证校验结果
type UserSession struct {
	CommonResp
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}
