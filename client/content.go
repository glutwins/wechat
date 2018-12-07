package client

// MsgSecCheck 内容检测
func (c *Client) MsgSecCheck(content string) *CommonResp {
	resp := &CommonResp{}
	if err := c.postJsonUrlFormat(map[string]string{"content": content}, resp, "https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s"); err != nil {
		resp.ErrCode = -1
		resp.ErrMsg = err.Error()
	}
	return resp
}
