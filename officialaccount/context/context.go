package context

import (
	"github.com/glutwins/wechat/credential"
	"github.com/glutwins/wechat/officialaccount/config"
)

// Context struct
type Context struct {
	*config.Config
	credential.AccessTokenHandle
}
