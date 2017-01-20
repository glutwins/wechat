package client

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"

	"github.com/glutwins/wechat/crypt"
)

type mchBaseRes struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

func (r mchBaseRes) Error() string {
	if r.ReturnCode != "SUCCESS" {
		return r.ReturnMsg
	}

	if r.ResultCode != "SUCCESS" {
		return r.ErrCode + " " + r.ErrCodeDes
	}

	return ""
}

type mchTransRes struct {
	mchBaseRes
	AppId      string `xml:"mch_appid"`
	MchId      string `xml:"mchid"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	TradeNo    string `xml:"partner_trade_no"`
	PayNo      string `xml:"payment_no"`
	PayTime    string `xml:"payment_time"`
}

type mchRepackRes struct {
	mchBaseRes
	Sign    string `xml:"sign"`
	TradeNo string `xml:"mch_billno"`
	MchId   string `xml:"mch_id"`
	AppId   string `xml:"wxappid"`
	OpenId  string `xml:"re_openid"`
	Amount  int    `xml:"total_amount"`
	PayNo   string `xml:"send_listid"`
}

func (c *Client) mchPost(url string, req map[string]interface{}, res error) error {
	sign := crypt.Md5Sign(req)
	req["sign"] = sign

	resp, err := c.mch.Post(url, "application/xml", bytes.NewReader([]byte(crypt.EncodeXml(req))))
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	if err = xml.Unmarshal(b, res); err != nil {
		return err
	}

	if res.Error() != "" {
		return res
	}

	return nil
}

func (c *Client) Transfer(openid string, tradeno string, name string, amount int, desc string) (string, error) {
	val := make(map[string]interface{})
	val["mch_appid"] = c.AppID
	val["mchid"] = c.MchId
	val["nonce_str"] = crypt.RandomStr(10)
	val["partner_trade_no"] = tradeno
	val["openid"] = openid
	val["check_name"] = "OPTION_CHECK"
	val["re_user_name"] = name
	val["amount"] = amount
	val["desc"] = desc
	val["spbill_create_ip"] = c.BillIp

	var res mchTransRes
	if err := c.mchPost(mchTransURL, val, &res); err != nil {
		return "", err
	}

	return res.PayNo, nil
}

func (c *Client) Redpack(openid, tradeno, mchname, wishing, act, remark, scene_id string, amount int, num int) (string, error) {
	val := make(map[string]interface{})
	val["nonce_str"] = crypt.RandomStr(10)
	val["mch_billno"] = tradeno
	val["mch_id"] = c.MchId
	val["wxappid"] = c.AppID
	val["send_name"] = mchname
	val["re_openid"] = openid
	val["total_amount"] = amount
	val["total_num"] = num
	val["wishing"] = wishing
	val["client_ip"] = c.BillIp
	val["act_name"] = act
	val["remark"] = remark

	if amount > 20000 {
		val["scene_id"] = scene_id
	}

	var res mchRepackRes
	if err := c.mchPost(mchRedpackURL, val, &res); err != nil {
		return "", nil
	}

	return res.PayNo, nil
}

func (c *Client) GroupRedpack(openid, tradeno, mchname, wishing, act, remark, scene_id string, amount int, num int) (string, error) {
	val := make(map[string]interface{})
	val["nonce_str"] = crypt.RandomStr(10)
	val["mch_billno"] = tradeno
	val["mch_id"] = c.MchId
	val["wxappid"] = c.AppID
	val["send_name"] = mchname
	val["re_openid"] = openid
	val["total_amount"] = amount
	val["total_num"] = num
	val["wishing"] = wishing
	val["client_ip"] = c.BillIp
	val["act_name"] = act
	val["remark"] = remark
	val["amt_type"] = "ALL_RAND"
	val["scene_id"] = scene_id

	var res mchRepackRes
	if err := c.mchPost(mchGroupRedpackURL, val, &res); err != nil {
		return "", nil
	}

	return res.PayNo, nil
}
