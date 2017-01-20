package context

import (
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/glutwins/wechat/crypt"
	"github.com/glutwins/wechat/message"
)

var (
	// ErrSign wechat sign error
	ErrSign = errors.New("invalid sign")
)

type Context interface {
	Query(key string) (string, bool)
	ParseRawMessage() (*message.MixMessage, error)
	Render([]byte)
	RenderString(string)
	RenderXML(interface{})
	RenderMessage(message.Reply)
}

func NewContext(req *http.Request, w http.ResponseWriter, cfg *crypt.WechatConfig) Context {
	ctx := &ctxImpl{Writer: w, Request: req, cfg: cfg}
	ctx.parse()
	return ctx
}

// Context struct
type ctxImpl struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	isSafeMode bool
	params     url.Values
	random     []byte
	cfg        *crypt.WechatConfig
}

func (ctx *ctxImpl) parse() {
	ctx.params = ctx.Request.URL.Query()
	if enctype, ok := ctx.Query("encrypt_type"); ok {
		ctx.isSafeMode = enctype == "aes"
	}
}

func (ctx *ctxImpl) ParseRawMessage() (*message.MixMessage, error) {
	timestamp, _ := ctx.Query("timestamp")
	nonce, _ := ctx.Query("nonce")
	signature, _ := ctx.Query("signature")

	if signature != crypt.Signature(ctx.cfg.Token, timestamp, nonce) {
		return nil, ErrSign
	}

	var rawXMLMsgBytes []byte
	var err error

	if ctx.isSafeMode {
		var encryptedXMLMsg message.EncryptedXMLMsg
		if err = xml.NewDecoder(ctx.Request.Body).Decode(&encryptedXMLMsg); err != nil {
			return nil, err
		}

		msgSignature, _ := ctx.Query("msg_signature")
		msgSignatureGen := crypt.Signature(ctx.cfg.Token, timestamp, nonce, encryptedXMLMsg.EncryptedMsg)
		if msgSignature != msgSignatureGen {
			return nil, ErrSign
		}

		//解密
		ctx.random, rawXMLMsgBytes, err = crypt.DecryptMsg(ctx.cfg.AppID, encryptedXMLMsg.EncryptedMsg, ctx.cfg.AESKey)
		if err != nil {
			return nil, err
		}
	} else {
		rawXMLMsgBytes, err = ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			return nil, err
		}
	}

	msg := &message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, msg)
	return msg, err
}

// GetQuery is like Query(), it returns the keyed url query value
func (ctx *ctxImpl) Query(key string) (string, bool) {
	if values, ok := ctx.params[key]; ok && len(values) > 0 {
		return values[0], true
	}
	return "", false
}

//Render render from bytes
func (ctx *ctxImpl) Render(bytes []byte) {
	ctx.Writer.WriteHeader(200)
	_, err := ctx.Writer.Write(bytes)
	if err != nil {
		panic(err)
	}
}

//String render from string
func (ctx *ctxImpl) RenderString(str string) {
	ctx.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.Render([]byte(str))
}

//XML render to xml
func (ctx *ctxImpl) RenderXML(obj interface{}) {
	ctx.Writer.Header().Set("Content-Type", "application/xml; charset=utf-8")
	bytes, err := xml.Marshal(obj)
	if err != nil {
		panic(err)
	}
	ctx.Render(bytes)
}

func (ctx *ctxImpl) RenderMessage(reply message.Reply) {
	var replyMsg interface{} = reply
	if ctx.isSafeMode {
		rawMsg, err := xml.Marshal(replyMsg)
		if err != nil {
			return
		}
		//安全模式下对消息进行加密
		encryptedMsg, err := crypt.EncryptMsg(ctx.random, rawMsg, ctx.cfg.AppID, ctx.cfg.AESKey)
		if err != nil {
			return
		}
		timestamp := time.Now().Unix()
		strts := strconv.FormatInt(timestamp, 64)
		nonce := crypt.RandomStr(16)
		msgSignature := crypt.Signature(ctx.cfg.Token, strts, nonce, string(encryptedMsg))
		replyMsg = message.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        nonce,
		}
	}
	ctx.RenderXML(replyMsg)
	return
}
