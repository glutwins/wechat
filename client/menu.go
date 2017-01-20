package client

//Button 菜单按钮
type Button struct {
	Type       string    `json:"type,omitempty"`
	Name       string    `json:"name,omitempty"`
	Key        string    `json:"key,omitempty"`
	URL        string    `json:"url,omitempty"`
	MediaID    string    `json:"media_id,omitempty"`
	SubButtons []*Button `json:"sub_button,omitempty"`
}

type Menu struct {
	pbtn *Button
	btns []*Button
}

func (menu *Menu) AddSub(name string) *Menu {
	if menu.pbtn != nil {
		return nil
	}

	sub := &Menu{pbtn: &Button{Name: name}}
	menu.btns = append(menu.btns, sub.pbtn)
	return sub
}

func (menu *Menu) addButton(btn *Button) {
	if menu.pbtn != nil {
		menu.pbtn.SubButtons = append(menu.pbtn.SubButtons, btn)
	} else {
		menu.btns = append(menu.btns, btn)
	}
}

// 添加click类型按钮
func (menu *Menu) AddClickButton(name, key string) {
	menu.addButton(&Button{Type: "click", Name: name, Key: key})
}

//添加view类型
func (menu *Menu) AddViewButton(name, url string) {
	menu.addButton(&Button{Type: "view", Name: name, URL: url})
}

// SetScanCodePushButton 扫码推事件
func (menu *Menu) AddScanCodePushButton(name, key string) {
	menu.addButton(&Button{Type: "scancode_push", Name: name, Key: key})
}

//SetScanCodeWaitMsgButton 设置 扫码推事件且弹出"消息接收中"提示框
func (menu *Menu) AddScanCodeWaitMsgButton(name, key string) {
	menu.addButton(&Button{Type: "scancode_waitmsg", Name: name, Key: key})
}

//SetPicSysPhotoButton 设置弹出系统拍照发图按钮
func (menu *Menu) AddPicSysPhotoButton(name, key string) {
	menu.addButton(&Button{Type: "pic_sysphoto", Name: name, Key: key})
}

//SetPicPhotoOrAlbumButton 设置弹出拍照或者相册发图类型按钮
func (menu *Menu) AddPicPhotoOrAlbumButton(name, key string) {
	menu.addButton(&Button{Type: "pic_photo_or_album", Name: name, Key: key})
}

// SetPicWeixinButton 设置弹出微信相册发图器类型按钮
func (menu *Menu) AddPicWeixinButton(name, key string) {
	menu.addButton(&Button{Type: "pic_weixin", Name: name, Key: key})
}

// SetLocationSelectButton 设置 弹出地理位置选择器 类型按钮
func (menu *Menu) AddLocationSelectButton(name, key string) {
	menu.addButton(&Button{Type: "location_select", Name: name, Key: key})
}

//SetMediaIDButton  设置 下发消息(除文本消息) 类型按钮
func (menu *Menu) AddMediaIDButton(name, mediaID string) {
	menu.addButton(&Button{Type: "media_id", Name: name, MediaID: mediaID})
}

//SetViewLimitedButton  设置 跳转图文消息URL 类型按钮
func (menu *Menu) AddViewLimitedButton(name, mediaID string) {
	menu.addButton(&Button{Type: "view_limited", Name: name, MediaID: mediaID})
}

//resMenuTryMatch 菜单匹配请求结果
type resMenuTryMatch struct {
	CommonResp

	Button []Button `json:"button"`
}

//ResMenu 查询菜单的返回数据
type ResMenu struct {
	CommonResp

	Menu struct {
		Button []Button `json:"button"`
		MenuID int64    `json:"menuid"`
	} `json:"menu"`
	Conditionalmenu []struct {
		Button    []Button  `json:"button"`
		MatchRule MatchRule `json:"matchrule"`
		MenuID    int64     `json:"menuid"`
	} `json:"conditionalmenu"`
}

//ResSelfMenuInfo 自定义菜单配置返回结果
type ResSelfMenuInfo struct {
	CommonResp

	IsMenuOpen   int32 `json:"is_menu_open"`
	SelfMenuInfo struct {
		Button []SelfMenuButton `json:"button"`
	} `json:"selfmenu_info"`
}

//SelfMenuButton 自定义菜单配置详情
type SelfMenuButton struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	URL       string `json:"url,omitempty"`
	Value     string `json:"value,omitempty"`
	SubButton struct {
		List []SelfMenuButton `json:"list"`
	} `json:"sub_button,omitempty"`
	NewsInfo struct {
		List []struct {
			Title      string `json:"title"`
			Author     string `json:"author"`
			Digest     string `json:"digest"`
			ShowCover  int32  `json:"show_cover"`
			CoverURL   string `json:"cover_url"`
			ContentURL string `json:"content_url"`
			SourceURL  string `json:"source_url"`
		} `json:"list"`
	} `json:"news_info,omitempty"`
}

//MatchRule 个性化菜单规则
type MatchRule struct {
	GroupID            int32  `json:"group_id,omitempty"`
	Sex                int32  `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType int32  `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}

//SetMenu 设置按钮
func (c *Client) SetMenu(menu *Menu) error {
	var commError CommonResp
	if err := c.postJsonUrlFormat(WachatReq{"button": menu.btns}, &commError, menuCreateURL); err != nil {
		return err
	}

	return commError.Error()
}

//GetMenu 获取菜单配置
func (c *Client) GetMenu() (*ResMenu, error) {
	resMenu := &ResMenu{}
	if err := c.getJsonUrlFormat(resMenu, menuGetURL); err != nil {
		return nil, err
	}
	return resMenu, resMenu.Error()
}

//DeleteMenu 删除菜单
func (c *Client) DeleteMenu() error {
	var commError CommonResp
	if err := c.getJsonUrlFormat(&commError, menuDeleteURL); err != nil {
		return err
	}
	return commError.Error()
}

//AddConditional 添加个性化菜单
func (c *Client) AddConditional(menu *Menu, matchRule *MatchRule) error {
	var commError CommonResp
	req := WachatReq{
		"button":    menu.btns,
		"matchrule": matchRule,
	}
	if err := c.postJsonUrlFormat(req, &commError, menuAddConditionalURL); err != nil {
		return err
	}
	return commError.Error()
}

//DeleteConditional 删除个性化菜单
func (c *Client) DeleteConditional(menuID int64) error {
	var commError CommonResp
	if err := c.postJsonUrlFormat(WachatReq{"menuid": menuID}, &commError, menuDeleteConditionalURL); err != nil {
		return err
	}
	return commError.Error()
}

//MenuTryMatch 菜单匹配
func (c *Client) MenuTryMatch(userID string) ([]Button, error) {
	var res resMenuTryMatch
	if err := c.postJsonUrlFormat(WachatReq{"user_id": userID}, &res, menuTryMatchURL); err != nil {
		return nil, err
	}

	return res.Button, res.Error()
}

//GetCurrentSelfMenuInfo 获取自定义菜单配置接口
func (c *Client) GetCurrentSelfMenuInfo() (*ResSelfMenuInfo, error) {
	resSelfMenuInfo := &ResSelfMenuInfo{}
	if err := c.getJsonUrlFormat(resSelfMenuInfo, menuSelfMenuInfoURL); err != nil {
		return nil, err
	}
	return resSelfMenuInfo, resSelfMenuInfo.Error()
}
