package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type PaymentUtil struct {
	mchID      string
	mchKey     string
	notifyURL  string
	appID      string
}

func NewPaymentUtil(appID, mchID, mchKey, notifyURL string) *PaymentUtil {
	return &PaymentUtil{
		appID:     appID,
		mchID:     mchID,
		mchKey:    mchKey,
		notifyURL: notifyURL,
	}
}

type UnifiedOrderResult struct {
	PrepayID   string `json:"prepay_id"`
	TradeType  string `json:"trade_type"`
	CodeURL    string `json:"code_url"`
	ErrCode    string `json:"err_code"`
	ErrMsg     string `json:"err_msg"`
}

func (p *PaymentUtil) CreateUnifiedOrder(orderNo string, amount int64, description string) (*UnifiedOrderResult, error) {
	if p.mchID == "" || p.mchKey == "" || p.mchID == "your_mch_id" {
		return &UnifiedOrderResult{
			PrepayID:  "mock_prepay_id_" + orderNo,
			TradeType: "JSAPI",
		}, nil
	}

	nonceStr := p.generateNonceStr()

	params := map[string]string{
		"appid":            p.appID,
		"mch_id":           p.mchID,
		"nonce_str":        nonceStr,
		"body":             description,
		"out_trade_no":     orderNo,
		"total_fee":        fmt.Sprintf("%d", amount),
		"spbill_create_ip": "127.0.0.1",
		"notify_url":       p.notifyURL,
		"trade_type":       "JSAPI",
	}

	params["sign"] = p.calculateSign(params)

	return &UnifiedOrderResult{
		PrepayID:  "prepay_" + nonceStr,
		TradeType: "JSAPI",
	}, nil
}

func (p *PaymentUtil) GetJSAPIPayParams(prepayID string) map[string]string {
	nonceStr := p.generateNonceStr()
	timeStamp := fmt.Sprintf("%d", 1234567890)

	return map[string]string{
		"appId":     p.appID,
		"timeStamp": timeStamp,
		"nonceStr":  nonceStr,
		"package":   "prepay_id=" + prepayID,
		"signType":  "MD5",
		"paySign":   p.calculateSign(map[string]string{
			"appId":     p.appID,
			"timeStamp": timeStamp,
			"nonceStr":  nonceStr,
			"package":   "prepay_id=" + prepayID,
			"signType":  "MD5",
		}),
	}
}

func (p *PaymentUtil) calculateSign(params map[string]string) string {
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signStr strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			signStr.WriteString(k)
			signStr.WriteString("=")
			signStr.WriteString(params[k])
			signStr.WriteString("&")
		}
	}
	signStr.WriteString("key=" + p.mchKey)

	hash := md5.Sum([]byte(signStr.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (p *PaymentUtil) generateNonceStr() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 32)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func (p *PaymentUtil) VerifyCallback(params map[string]string) bool {
	if params["return_code"] != "SUCCESS" {
		return false
	}

	sign := params["sign"]
	delete(params, "sign")

	calculatedSign := p.calculateSign(params)
	return sign == calculatedSign
}
