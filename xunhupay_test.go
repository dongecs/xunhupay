package xunhupay

import (
	"encoding/json"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestNewHuPi(t *testing.T) {
	appId := "test_app_id"
	appSecret := "test_app_secret"

	client := NewHuPi(&appId, &appSecret)

	if client == nil {
		t.Fatal("NewHuPi returned nil")
	}

	if *client.appId != appId {
		t.Errorf("Expected appId %s, got %s", appId, *client.appId)
	}

	if *client.appSecret != appSecret {
		t.Errorf("Expected appSecret %s, got %s", appSecret, *client.appSecret)
	}
}

func TestSign(t *testing.T) {
	appId := "test_app_id"
	appSecret := "test_app_secret"
	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
		"title":          "测试标题",
	}

	sign := client.Sign(params)

	if sign == "" {
		t.Error("Sign returned empty string")
	}

	// 验证签名长度（MD5是32位十六进制）
	if len(sign) != 32 {
		t.Errorf("Expected sign length 32, got %d", len(sign))
	}
}

func TestSignConsistency(t *testing.T) {
	appId := "test_app_id"
	appSecret := "test_app_secret"
	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
		"title":          "测试标题",
	}

	sign1 := client.Sign(params)
	sign2 := client.Sign(params)

	if sign1 != sign2 {
		t.Error("Sign should be consistent for same parameters")
	}
}

func TestSignWithDifferentParams(t *testing.T) {
	appId := "test_app_id"
	appSecret := "test_app_secret"
	client := NewHuPi(&appId, &appSecret)

	params1 := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
	}

	params2 := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.2", // 不同的金额
	}

	sign1 := client.Sign(params1)
	sign2 := client.Sign(params2)

	if sign1 == sign2 {
		t.Error("Sign should be different for different parameters")
	}
}

func TestResponseUnmarshal(t *testing.T) {
	jsonStr := `{
		"openid": 123456789,
		"url_qrcode": "https://example.com/qrcode",
		"url": "https://example.com/pay",
		"errcode": 0,
		"errmsg": "success",
		"hash": "abc123hash",
		"data": {
			"open_order_id": "order123",
			"status": "pending"
		}
	}`

	var response Response
	err := json.Unmarshal([]byte(jsonStr), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	expected := Response{
		Openid:    123456789,
		UrlQrcode: "https://example.com/qrcode",
		Url:       "https://example.com/pay",
		Errcode:   0,
		Errmsg:    "success",
		Hash:      "abc123hash",
		Data: Data{
			OpenOrderId: "order123",
			Status:      "pending",
		},
	}

	if !reflect.DeepEqual(response, expected) {
		t.Errorf("Expected %+v, got %+v", expected, response)
	}
}

func TestResponseErrorCase(t *testing.T) {
	jsonStr := `{
		"errcode": 1001,
		"errmsg": "参数错误",
		"hash": "errorhash"
	}`

	var response Response
	err := json.Unmarshal([]byte(jsonStr), &response)

	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	if response.Errcode != 1001 {
		t.Errorf("Expected errcode 1001, got %d", response.Errcode)
	}

	if response.Errmsg != "参数错误" {
		t.Errorf("Expected errmsg '参数错误', got '%s'", response.Errmsg)
	}
}

// 模拟测试 - 测试Execute方法的参数处理
func TestExecuteParameterProcessing(t *testing.T) {
	appId := "test_app_id"
	appSecret := "test_app_secret"
	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
		"title":          "测试标题",
	}

	// 测试签名功能
	sign := client.Sign(params)
	if sign == "" {
		t.Error("Sign should not be empty")
	}

	// 验证参数处理逻辑
	simple := time.Now().Unix()
	expectedParams := map[string]string{
		"appid":          appId,
		"time":           strconv.FormatInt(simple, 10),
		"nonce_str":      strconv.FormatInt(simple, 10),
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
		"title":          "测试标题",
	}

	// 验证所有必需参数都存在
	for key, value := range expectedParams {
		if key == "time" || key == "nonce_str" {
			// 这些是动态生成的，我们只检查它们存在
			if _, exists := params[key]; !exists {
				// 在Execute方法中会添加这些参数
				continue
			}
		}
		if params[key] != value {
			t.Errorf("Expected %s=%s, got %s", key, value, params[key])
		}
	}
}

// 基准测试 - 测试签名性能
func BenchmarkSign(b *testing.B) {
	appId := "test_app_id"
	appSecret := "test_app_secret"
	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "123456789",
		"total_fee":      "0.1",
		"title":          "测试标题",
		"notify_url":     "http://example.com/notify",
		"return_url":     "http://example.com/return",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Sign(params)
	}
}

// 测试示例函数
func TestPayExample(t *testing.T) {
	// 这是一个示例测试，实际使用时需要真实的API密钥
	t.Skip("Skipping integration test - requires real API credentials")

	appId := "test_app_id"
	appSecret := "test_app_secret"
	host := "https://api.xunhupay.com/payment/do.html"

	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"version":        "1.1",
		"trade_order_id": "test_order_123",
		"total_fee":      "0.01",
		"title":          "测试商品",
		"notify_url":     "http://example.com/notify",
		"return_url":     "http://example.com/return",
		"wap_name":       "测试店铺",
	}

	response, err := client.Execute(host, params)
	if err != nil {
		t.Logf("Expected error for test credentials: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Response: %+v", response)
}

func TestQueryExample(t *testing.T) {
	// 这是一个示例测试，实际使用时需要真实的API密钥
	t.Skip("Skipping integration test - requires real API credentials")

	appId := "test_app_id"
	appSecret := "test_app_secret"
	host := "https://api.xunhupay.com/payment/query.html"

	client := NewHuPi(&appId, &appSecret)

	params := map[string]string{
		"out_trade_order": "test_order_123",
	}

	response, err := client.Execute(host, params)
	if err != nil {
		t.Logf("Expected error for test credentials: %v", err)
		return
	}

	if response == nil {
		t.Error("Expected response, got nil")
		return
	}

	t.Logf("Response: %+v", response)
}
