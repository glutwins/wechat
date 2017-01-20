package server

import (
	"net/http"
	"time"

	"github.com/glutwins/wechat/context"
	"github.com/glutwins/wechat/crypt"
	"github.com/glutwins/wechat/message"
)

type Handler func(*message.MixMessage) message.Reply

//Server struct
type Server struct {
	*crypt.WechatConfig
	messageHandler Handler
}

//NewServer init
func NewServer(cfg *crypt.WechatConfig, h Handler) *Server {
	srv := new(Server)
	srv.WechatConfig = cfg
	srv.messageHandler = h
	return srv
}

//Serve 处理微信的请求消息
func (srv *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := context.NewContext(req, w, srv.WechatConfig)
	echostr, exists := ctx.Query("echostr")
	if exists {
		ctx.RenderString(echostr)
		return
	}

	msg, err := ctx.ParseRawMessage()
	if err != nil {
		panic(err)
	}

	reply := srv.messageHandler(msg)
	if reply != nil {
		reply.SetToUserName(msg.FromUserName)
		reply.SetFromUserName(msg.ToUserName)
		reply.SetCreateTime(time.Now().Unix())
		ctx.RenderMessage(reply)
	}
}
