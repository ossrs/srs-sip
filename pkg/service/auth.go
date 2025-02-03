package service

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
)

// AuthInfo 存储解析后的认证信息
type AuthInfo struct {
	Username  string
	Realm     string
	Nonce     string
	URI       string
	Response  string
	Algorithm string
	Method    string
}

// GenerateNonce 生成随机 nonce 字符串
func GenerateNonce() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}

// ParseAuthorization 解析 SIP Authorization 头
// Authorization: Digest username="34020000001320000001",realm="3402000000",
// nonce="44010b73623249f6916a6acf7c316b8e",uri="sip:34020000002000000001@3402000000",
// response="e4ca3fdc5869fa1c544ea7af60014444",algorithm=MD5
func ParseAuthorization(auth string) *AuthInfo {
	auth = strings.TrimPrefix(auth, "Digest ")
	parts := strings.Split(auth, ",")
	result := &AuthInfo{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if !strings.Contains(part, "=") {
			continue
		}
		
		kv := strings.SplitN(part, "=", 2)
		key := strings.TrimSpace(kv[0])
		value := strings.Trim(strings.TrimSpace(kv[1]), "\"")
		
		switch key {
		case "username":
			result.Username = value
		case "realm":
			result.Realm = value
		case "nonce":
			result.Nonce = value
		case "uri":
			result.URI = value
		case "response":
			result.Response = value
		case "algorithm":
			result.Algorithm = value
		}
	}
	
	return result
}

// ValidateAuth 验证 SIP 认证信息
func ValidateAuth(authInfo *AuthInfo, password string) bool {
	if authInfo == nil {
		return false
	}

	// 默认方法为 REGISTER
	method := "REGISTER"
	if authInfo.Method != "" {
		method = authInfo.Method
	}

	// 计算 MD5 哈希
	ha1 := md5Hex(authInfo.Username + ":" + authInfo.Realm + ":" + password)
	ha2 := md5Hex(method + ":" + authInfo.URI)
	correctResponse := md5Hex(ha1 + ":" + authInfo.Nonce + ":" + ha2)
	
	return authInfo.Response == correctResponse
}

// md5Hex 计算字符串的 MD5 哈希值并返回十六进制字符串
func md5Hex(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
} 