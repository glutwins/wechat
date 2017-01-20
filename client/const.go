package client

const (
	accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	getTicketURL   = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"

	addNewsURL     = "https://api.weixin.qq.com/cgi-bin/material/add_news?access_token=%s"
	addMaterialURL = "https://api.weixin.qq.com/cgi-bin/material/add_material?access_token=%s&type=%s"
	delMaterialURL = "https://api.weixin.qq.com/cgi-bin/material/del_material?access_token=%s"

	mediaUploadURL      = "https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	mediaUploadImageURL = "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s"
	mediaGetURL         = "https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"

	menuCreateURL            = "https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s"
	menuGetURL               = "https://api.weixin.qq.com/cgi-bin/menu/get?access_token=%s"
	menuDeleteURL            = "https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s"
	menuAddConditionalURL    = "https://api.weixin.qq.com/cgi-bin/menu/addconditional?access_token=%s"
	menuDeleteConditionalURL = "https://api.weixin.qq.com/cgi-bin/menu/delconditional?access_token=%s"
	menuTryMatchURL          = "https://api.weixin.qq.com/cgi-bin/menu/trymatch?access_token=%s"
	menuSelfMenuInfoURL      = "https://api.weixin.qq.com/cgi-bin/get_current_selfmenu_info?access_token=%s"

	redirectOauthURL      = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=%s&state=%s#wechat_redirect"
	userAccessTokenURL    = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	refreshAccessTokenURL = "https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s"
	userInfoURL           = "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN"
	checkAccessTokenURL   = "https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s"

	qrcodeURL = "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s"

	mchTransURL        = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
	mchRedpackURL      = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
	mchGroupRedpackURL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendgroupredpack"
)

//MediaType 媒体文件类型
type MediaType string

const (
	//MediaTypeImage 媒体文件:图片
	MediaTypeImage MediaType = "image"
	//MediaTypeVoice 媒体文件:声音
	MediaTypeVoice = "voice"
	//MediaTypeVideo 媒体文件:视频
	MediaTypeVideo = "video"
	//MediaTypeThumb 媒体文件:缩略图
	MediaTypeThumb = "thumb"
)

const (
	kQrIntFormat      = `{"expire_seconds": %d, "action_name": "QR_SCENE", "action_info": {"scene": {"scene_id": %d}}}`
	kQrIntLimitFormat = `{"action_name": "QR_LIMIT_SCENE", "action_info": {"scene": {"scene_id": %d}}}`
	kQrStrLimitFormat = `{"action_name": "QR_LIMIT_STR_SCENE", "action_info": {"scene": {"scene_str": "%s"}}}`
)
