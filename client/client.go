package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/glutwins/webclient"
	"github.com/glutwins/wechat/cache"
	"github.com/glutwins/wechat/crypt"
)

type WachatReq map[string]interface{}

type Client struct {
	*crypt.WechatConfig
	data cache.Cache
	mch  *http.Client
}

func NewClient(cfg *crypt.WechatConfig, d cache.Cache) *Client {
	c := &Client{WechatConfig: cfg, data: d}
	if cfg.MchId != "" {
		cert, err := tls.LoadX509KeyPair(cfg.ServerPem, cfg.ServerKey)
		if err != nil {
			panic(err)
		}
		certBytes, err := ioutil.ReadFile(cfg.ClientCert)
		if err != nil {
			panic(err)
		}
		clientCertPool := x509.NewCertPool()
		if !clientCertPool.AppendCertsFromPEM(certBytes) {
			panic("failed to parse root certificate")
		}
		conf := &tls.Config{
			RootCAs:            clientCertPool,
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		}
		c.mch = &http.Client{Transport: &http.Transport{TLSClientConfig: conf}}
	}
	return c
}

type JsConfig struct {
	AppID     string
	TimeStamp int64
	NonceStr  string
	Signature string
}

type cacheClient struct {
	Token  accessToken
	Ticket resTicket
}

//GetAccessToken 获取access_token
func (c *Client) GetAccessToken() (string, error) {
	var token accessToken
	if err := c.data.Get("token", &token); err == nil {
		if token.valid() {
			return token.AccessToken, nil
		}
	}

	c.data.Lock("token")
	defer c.data.Unlock("token")

	url := fmt.Sprintf(accessTokenURL, c.AppID, c.AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GetAccessToken httpcode=%d", resp.StatusCode)
	}

	if b, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else if err = json.Unmarshal(b, &token); err != nil {
		return "", err
	}
	token.LastUpTs = time.Now().Unix()
	c.data.Set("token", token)

	return token.AccessToken, token.Error()
}

func (c *Client) formatUrlWithAccessToken(base string, args ...interface{}) (string, error) {
	if token, err := c.GetAccessToken(); err != nil {
		return "", err
	} else {
		return fmt.Sprintf(base, append([]interface{}{token}, args...)...), nil
	}
}

func (c *Client) postJsonUrlFormat(req interface{}, res interface{}, url string, args ...interface{}) error {
	if uri, err := c.formatUrlWithAccessToken(url, args...); err != nil {
		return err
	} else if response, err := c.jsonPost(uri, req); err != nil {
		return err
	} else if err := json.Unmarshal(response, res); err != nil {
		return err
	}
	return nil
}

func (c *Client) jsonPost(uri string, obj interface{}) ([]byte, error) {
	var jsonData []byte
	var err error

	if str, ok := obj.(string); ok {
		jsonData = []byte(str)
	} else {
		jsonData, err = json.Marshal(obj)
		if err != nil {
			return nil, err
		}
	}
	return c.doPost(uri, "application/json;charset=utf-8", bytes.NewBuffer(jsonData))
}

func (c *Client) formPost(uri string, data map[string]string) ([]byte, error) {
	body := bytes.NewBuffer(nil)
	for k, v := range data {
		body.WriteString(k)
		body.WriteString("=")
		body.WriteString(url.QueryEscape(v))
		body.WriteString("&")
	}
	body.Truncate(body.Len() - 1)
	return c.doPost(uri, "application/x-www-form-urlencoded;charset=utf-8", body)
}

func (c *Client) doPost(uri, contentType string, r io.Reader) ([]byte, error) {
	response, err := http.Post(uri, contentType, r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http get error : uri=%v , statusCode=%v", uri, response.StatusCode)
	}
	return ioutil.ReadAll(response.Body)
}

func (c *Client) getJsonUrlFormat(res interface{}, url string, args ...interface{}) error {
	uri, err := c.formatUrlWithAccessToken(url, args...)
	if err != nil {
		return err
	}

	if b, err := webclient.DoGet(uri); err != nil {
		return err
	} else {
		return json.Unmarshal(b, res)
	}
}

//uri 为当前网页地址
func (js *Client) GetJsConfig(uri string) (config *JsConfig, err error) {
	config = new(JsConfig)
	var ticketStr string
	ticketStr, err = js.getTicket()
	if err != nil {
		return
	}

	nonceStr := crypt.RandomStr(16)
	timestamp := time.Now().Unix()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticketStr, nonceStr, timestamp, uri)
	sigStr := crypt.Signature(str)

	config.AppID = js.AppID
	config.NonceStr = nonceStr
	config.TimeStamp = timestamp
	config.Signature = sigStr
	return
}

//getTicket 获取jsapi_tocket全局缓存
func (c *Client) getTicket() (string, error) {
	var ticket resTicket
	if err := c.data.Get("ticket", &ticket); err == nil {
		if ticket.valid() {
			return ticket.Ticket, nil
		}
	}

	c.data.Lock("ticket")
	defer c.data.Unlock("ticket")

	if err := c.getJsonUrlFormat(&ticket, getTicketURL); err != nil {
		return "", err
	}

	ticket.LastUpTs = time.Now().Unix()
	c.data.Set("ticket", ticket)

	return ticket.Ticket, ticket.Error()
}

func (c *Client) NewQrCode(sceneId interface{}, expire int) (*QrCode, error) {
	var msg string
	if _, ok := sceneId.(string); ok {
		if expire == 0 {
			msg = fmt.Sprintf(kQrStrLimitFormat, sceneId)
		} else {
			return nil, fmt.Errorf("not support expire str qrcode")
		}
	} else {
		if expire == 0 {
			msg = fmt.Sprintf(kQrIntLimitFormat, sceneId)
		} else {
			msg = fmt.Sprintf(kQrIntFormat, expire, sceneId)
		}
	}

	var qrcode = &QrCode{}
	if err := c.postJsonUrlFormat(msg, qrcode, qrcodeURL); err != nil {
		return nil, err
	}
	return qrcode, qrcode.Error()
}
