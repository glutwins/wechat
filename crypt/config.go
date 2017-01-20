package crypt

type WechatConfig struct {
	AppID     string
	AppSecret string
	Token     string
	AESKey    string

	// MchId 商户ID，用于支付
	MchId string
	// ServerPem 微信服务器证书
	ServerPem string
	// ServerKey 微信服务器私钥
	ServerKey string
	// ClientCert 用于支付的商户证书
	ClientCert string
	BillIp     string
}
