package crypt

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"
	"time"
)

//Signature sha1签名
func Signature(params ...string) string {
	sort.Strings(params)
	h := sha1.New()
	for _, s := range params {
		io.WriteString(h, s)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5Sign(val map[string]interface{}) string {
	strs := make([]string, 0, len(val))
	for k, v := range val {
		strs = append(strs, fmt.Sprintf("%s=%v", k, v))
	}

	h := md5.New()
	h.Write([]byte(strings.Join(strs, "&")))
	return hex.EncodeToString(h.Sum(nil))
}

func EncodeXml(val map[string]interface{}) string {
	var buf bytes.Buffer
	buf.WriteString("<xml>")
	for k, v := range val {
		buf.WriteString("<" + k + ">")
		buf.WriteString(fmt.Sprint(v))
		buf.WriteString("</" + k + ">")
	}
	buf.WriteString("</xml>")
	return buf.String()
}

func RandomStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
