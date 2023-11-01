package jwt

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gitee.com/Kashimura/go-baka-control/services/tools"
	"github.com/google/uuid"
)

const JWT_KEY = "Kashimura"

// Jwt 数据签名加密
func signature(data string, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return tools.Base64UrlEncode(h.Sum(nil))
}

// 获取数据 Token
func GetToken(payload *PayLoad, timeout int64) string {
	nowtime := time.Now().Unix()
	// 更新 签名时间
	payload.ISS = JWT_KEY
	payload.IAT = nowtime
	payload.EXP = nowtime + timeout*60
	payload.NBF = nowtime
	payload.JTI = "JWT" + uuid.New().String() + strconv.FormatInt(nowtime, 10)
	// 数据打包
	payloaddata, _ := json.Marshal(payload)
	base64payload := tools.Base64UrlEncode(payloaddata)
	return base64payload + "." + signature(base64payload, JWT_KEY)
}

// 根据文件字节流获取
func GetFileToken(fileid int64) string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, fileid)
	longBytes := buf.Bytes()
	base64File := base64.StdEncoding.EncodeToString(longBytes)
	nowtime := time.Now().Unix()
	endtime := nowtime + 21600 // 6个小时
	info := fmt.Sprintf("%s.%d.%d", base64File, nowtime, endtime)
	sign := signature(info, JWT_KEY)
	return tools.Base64UrlEncode([]byte(info + "." + sign))
}

// 校验 Token 是否有效
func VerifyToken(token string) (*PayLoad, error) {
	// 检查Token
	if token == "" {
		return nil, errors.New("jwtz is null")
	}
	// 是否为有效token
	tokens := strings.Split(token, ".")
	if len(tokens) != 2 {
		return nil, errors.New("token format error")
	}
	// 获取token信息
	base64payload := tokens[0]
	sign := tokens[1]
	// 校验签名
	if signature(base64payload, JWT_KEY) != sign {
		return nil, errors.New("not a valid token")
	}
	// 获取信息
	var payload PayLoad
	if err := json.Unmarshal(tools.Base64UrlDecode(base64payload), &payload); err != nil {
		return nil, errors.New("token parsing error")
	}
	// 获取当前服务器时间
	nowtime := time.Now().Unix()
	// 签发时间大于当前服务器时间验证失败
	if payload.IAT != 0 && payload.IAT > nowtime {
		return nil, errors.New("iat timeout")
	}
	// 过期时间小于当前服务器时间验证失败
	if payload.EXP != 0 && payload.EXP < nowtime {
		return nil, errors.New("exp timeout")
	}
	// 该nbf时间之前不接收处理该Token
	if payload.NBF != 0 && payload.NBF > nowtime {
		return nil, errors.New("nbf timeout")
	}
	return &payload, nil
}
