package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"fmt"
)

//EncryptMsg 加密消息
func EncryptMsg(random, rawXMLMsg []byte, appID, aesKey string) (encrtptMsg []byte, err error) {
	var key []byte
	key, err = aesKeyDecode(aesKey)
	if err != nil {
		return nil, err
	}
	ciphertext := AESEncryptMsg(random, rawXMLMsg, appID, key)
	encrtptMsg = []byte(base64.StdEncoding.EncodeToString(ciphertext))
	return
}

//AESEncryptMsg ciphertext = AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId]
//参考：github.com/chanxuehong/wechat.v2
func AESEncryptMsg(random, rawXMLMsg []byte, appID string, aesKey []byte) (ciphertext []byte) {
	const (
		BlockSize = 32            // PKCS#7
		BlockMask = BlockSize - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
	)

	appIDOffset := 20 + len(rawXMLMsg)
	contentLen := appIDOffset + len(appID)
	amountToPad := BlockSize - contentLen&BlockMask
	plaintextLen := contentLen + amountToPad

	plaintext := make([]byte, plaintextLen)

	// 拼接
	copy(plaintext[:16], random)
	binary.BigEndian.PutUint32(plaintext[16:20], uint32(len(rawXMLMsg)))
	copy(plaintext[20:], rawXMLMsg)
	copy(plaintext[appIDOffset:], appID)

	// PKCS#7 补位
	for i := contentLen; i < plaintextLen; i++ {
		plaintext[i] = byte(amountToPad)
	}

	// 加密
	block, err := aes.NewCipher(aesKey[:])
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, plaintext)

	ciphertext = plaintext
	return
}

//DecryptMsg 消息解密
func DecryptMsg(appID, encryptedMsg, aesKey string) (random, rawMsgXMLBytes []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic error: err=%v", e)
			return
		}
	}()
	var encryptedMsgBytes, key, getAppIDBytes []byte
	encryptedMsgBytes, err = base64.StdEncoding.DecodeString(encryptedMsg)
	if err != nil {
		return
	}
	key, err = aesKeyDecode(aesKey)
	if err != nil {
		panic(err)
	}
	random, rawMsgXMLBytes, getAppIDBytes, err = AESDecryptMsg(encryptedMsgBytes, key)
	if err != nil {
		err = fmt.Errorf("消息解密失败,%v", err)
		return
	}
	if appID != string(getAppIDBytes) {
		err = fmt.Errorf("消息解密校验APPID失败")
		return
	}
	return
}

func aesKeyDecode(encodedAESKey string) (key []byte, err error) {
	if len(encodedAESKey) != 43 {
		err = fmt.Errorf("the length of encodedAESKey must be equal to 43")
		return
	}
	key, err = base64.StdEncoding.DecodeString(encodedAESKey + "=")
	if err != nil {
		return
	}
	if len(key) != 32 {
		err = fmt.Errorf("encodingAESKey invalid")
		return
	}
	return
}

// AESDecryptMsg ciphertext = AES_Encrypt[random(16B) + msg_len(4B) + rawXMLMsg + appId]
//参考：github.com/chanxuehong/wechat.v2
func AESDecryptMsg(ciphertext []byte, aesKey []byte) (random, rawXMLMsg, appID []byte, err error) {
	const (
		BlockSize = 32            // PKCS#7
		BlockMask = BlockSize - 1 // BLOCK_SIZE 为 2^n 时, 可以用 mask 获取针对 BLOCK_SIZE 的余数
	)

	if len(ciphertext) < BlockSize {
		err = fmt.Errorf("the length of ciphertext too short: %d", len(ciphertext))
		return
	}
	if len(ciphertext)&BlockMask != 0 {
		err = fmt.Errorf("ciphertext is not a multiple of the block size, the length is %d", len(ciphertext))
		return
	}

	plaintext := make([]byte, len(ciphertext)) // len(plaintext) >= BLOCK_SIZE

	// 解密
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCDecrypter(block, aesKey[:16])
	mode.CryptBlocks(plaintext, ciphertext)

	// PKCS#7 去除补位
	amountToPad := int(plaintext[len(plaintext)-1])
	if amountToPad < 1 || amountToPad > BlockSize {
		err = fmt.Errorf("the amount to pad is incorrect: %d", amountToPad)
		return
	}
	plaintext = plaintext[:len(plaintext)-amountToPad]

	// 反拼接
	if len(plaintext) <= 20 {
		err = fmt.Errorf("plaintext too short, the length is %d", len(plaintext))
		return
	}
	rawXMLMsgLen := int(binary.BigEndian.Uint32(plaintext[16:20]))
	if rawXMLMsgLen < 0 {
		err = fmt.Errorf("incorrect msg length: %d", rawXMLMsgLen)
		return
	}
	appIDOffset := 20 + rawXMLMsgLen
	if len(plaintext) <= appIDOffset {
		err = fmt.Errorf("msg length too large: %d", rawXMLMsgLen)
		return
	}

	random = plaintext[:16:20]
	rawXMLMsg = plaintext[20:appIDOffset:appIDOffset]
	appID = plaintext[appIDOffset:]
	return
}
