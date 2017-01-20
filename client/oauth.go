package client

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/glutwins/webclient"
)

//GetRedirectURL 获取跳转的url地址
func (c *Client) GetRedirectURL(redirectURI, scope, state string) string {
	//url encode
	urlStr := url.QueryEscape(redirectURI)
	return fmt.Sprintf(redirectOauthURL, c.AppID, urlStr, scope, state)
}

// GetUserAccessToken 通过网页授权的code 换取access_token(区别于context中的access_token)
func (c *Client) GetUserAccessToken(code string) (*UserAccessToken, error) {
	if b, err := webclient.DoGet(fmt.Sprintf(userAccessTokenURL, c.AppID, c.AppSecret, code)); err != nil {
		return nil, err
	} else {
		token := &UserAccessToken{}
		if err = json.Unmarshal(b, token); err != nil {
			return nil, err
		}
		return token, token.Error()
	}
}

//RefreshAccessToken 刷新access_token
func (c *Client) RefreshAccessToken(refreshToken string) (*UserAccessToken, error) {
	if b, err := webclient.DoGet(fmt.Sprintf(refreshAccessTokenURL, c.AppID, refreshToken)); err != nil {
		return nil, err
	} else {
		token := &UserAccessToken{}
		if err = json.Unmarshal(b, token); err != nil {
			return nil, err
		}
		return token, token.Error()
	}
}

//CheckAccessToken 检验access_token是否有效
func (c *Client) CheckAccessToken(accessToken, openID string) error {
	var result CommonResp
	if err := c.getJsonUrlFormat(&result, checkAccessTokenURL, openID); err != nil {
		return err
	}
	return result.Error()
}

//GetUserInfo 如果scope为 snsapi_userinfo 则可以通过此方法获取到用户基本信息
func (c *Client) GetUserInfo(accessToken, openID string) (*UserInfo, error) {
	user := &UserInfo{}
	uri := fmt.Sprintf(userInfoURL, accessToken, openID)

	if b, err := webclient.DoGet(uri); err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &user); err != nil {
		return nil, err
	}
	return user, user.Error()
}
