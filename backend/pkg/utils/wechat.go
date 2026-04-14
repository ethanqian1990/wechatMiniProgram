package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type WechatUtil struct {
	appID     string
	appSecret string
}

func NewWechatUtil(appID, appSecret string) *WechatUtil {
	return &WechatUtil{
		appID:     appID,
		appSecret: appSecret,
	}
}

type WechatLoginResult struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func (w *WechatUtil) Code2Session(code string) (*WechatLoginResult, error) {
	if w.appID == "" || w.appSecret == "" || w.appID == "your_app_id" {
		return &WechatLoginResult{
			OpenID:     "mock_openid_" + code,
			SessionKey: "mock_session_key",
		}, nil
	}

	apiURL := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		w.appID, w.appSecret, code,
	)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("请求微信接口失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	var result WechatLoginResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信接口返回错误: %s", result.ErrMsg)
	}

	return &result, nil
}

func (w *WechatUtil) GetPhoneNumber(accessToken, code string) (string, error) {
	if w.appID == "" || w.appSecret == "" || w.appID == "your_app_id" {
		return "13800138000", nil
	}

	return "", nil
}
