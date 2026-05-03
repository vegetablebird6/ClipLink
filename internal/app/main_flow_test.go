package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xiaojiu/cliplink/internal/config"
	"github.com/xiaojiu/cliplink/internal/infra/db"
)

type testAPIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Success bool            `json:"success"`
}

func TestMainAPIClipboardFlow(t *testing.T) {
	router := newMainFlowRouter(t)
	channelID := "main-flow-channel"
	deviceID := "main-flow-device"

	postJSON(t, router, "/api/channel", "", http.StatusOK, map[string]any{
		"channel_id": channelID,
	})

	doJSON(t, router, http.MethodPost, "/api/channel/verify", "", http.StatusOK, map[string]any{
		"channel_id": channelID,
	})

	device := postJSON(t, router, "/api/devices", channelID, http.StatusOK, map[string]any{
		"device_id":   deviceID,
		"device_name": "Windows Chrome",
		"device_type": "desktop",
	})
	assertDataField(t, device, "id", deviceID)

	textItem := saveClipboardViaAPI(t, router, channelID, deviceID, "Alpha plain text", "text", "Alpha")
	saveClipboardViaAPI(t, router, channelID, deviceID, "https://example.com", "link", "Example")
	codeItem := saveClipboardViaAPI(t, router, channelID, deviceID, "fmt.Println(\"hi\")", "code", "Snippet")

	current := getJSON(t, router, "/api/clipboard/current", channelID, http.StatusOK)
	assertDataField(t, dataObject(t, current), "id", stringField(t, codeItem, "id"))

	firstHistory := getJSON(t, router, "/api/clipboard/history?size=2", channelID, http.StatusOK)
	assertDataBool(t, firstHistory, "has_more", true)
	firstItems := dataItems(t, firstHistory)
	assertLen(t, firstItems, 2, "first history page")
	assertField(t, firstItems[0], "id", stringField(t, codeItem, "id"))

	after := stringFieldFromMap(t, firstItems[1], "created_at")
	afterID := stringFieldFromMap(t, firstItems[1], "id")
	secondHistory := getJSON(t, router, fmt.Sprintf("/api/clipboard/history?size=2&after=%s&after_id=%s", url.QueryEscape(after), url.QueryEscape(afterID)), channelID, http.StatusOK)
	assertDataBool(t, secondHistory, "has_more", false)
	assertLen(t, dataItems(t, secondHistory), 1, "second history page")

	textHistory := getJSON(t, router, "/api/clipboard/type/text?size=2", channelID, http.StatusOK)
	textItems := dataItems(t, textHistory)
	assertLen(t, textItems, 1, "text history page")
	assertField(t, textItems[0], "id", stringField(t, textItem, "id"))

	search := getJSON(t, router, "/api/clipboard/search?q=Alpha&page=1&size=10", channelID, http.StatusOK)
	searchItems := dataItems(t, search)
	assertLen(t, searchItems, 1, "search results")
	assertField(t, searchItems[0], "id", stringField(t, textItem, "id"))

	putJSON(t, router, "/api/clipboard/"+stringField(t, textItem, "id")+"/favorite", channelID, http.StatusOK, map[string]any{
		"favorite":  true,
		"device_id": deviceID,
	})

	favorites := getJSON(t, router, "/api/clipboard/favorites", channelID, http.StatusOK)
	favoriteItems := dataArray(t, favorites)
	assertLen(t, favoriteItems, 1, "favorites")
	assertField(t, favoriteItems[0], "id", stringField(t, textItem, "id"))

	putJSON(t, router, "/api/clipboard/"+stringField(t, textItem, "id")+"/favorite", channelID, http.StatusOK, map[string]any{
		"favorite":  false,
		"device_id": deviceID,
	})
	emptyFavorites := getJSON(t, router, "/api/clipboard/favorites", channelID, http.StatusOK)
	assertLen(t, dataArray(t, emptyFavorites), 0, "favorites after unfavorite")

	updated := putJSON(t, router, "/api/clipboard/"+stringField(t, textItem, "id"), channelID, http.StatusOK, map[string]any{
		"title":          "Alpha Updated",
		"content":        "Alpha updated body",
		"type":           "text",
		"device_id":      deviceID,
		"device_type":    "desktop",
		"content_format": "plain",
	})
	assertDataField(t, updated, "content", "Alpha updated body")

	// 验证部分更新：只改 title 不清空其他字段
	titleOnly := putJSON(t, router, "/api/clipboard/"+stringField(t, textItem, "id"), channelID, http.StatusOK, map[string]any{
		"title":     "Title Only Update",
		"device_id": deviceID,
	})
	assertDataField(t, titleOnly, "title", "Title Only Update")
	assertDataField(t, titleOnly, "content", "Alpha updated body")   // content 保持不变
	assertDataField(t, titleOnly, "type", "text")                   // type 保持不变

	cleared := putJSON(t, router, "/api/clipboard/"+stringField(t, textItem, "id"), channelID, http.StatusOK, map[string]any{
		"content_html":   "",
		"content_format": "plain",
		"device_id":      deviceID,
	})
	// omitempty 导致空 content_html 不出现在响应中（与不设置效果一致）
	if html, exists := cleared["content_html"]; exists {
		if s, ok := html.(string); !ok || s != "" {
			t.Fatalf("expected empty content_html, got %v", html)
		}
	}

	putJSON(t, router, "/api/devices/"+deviceID+"/name", channelID, http.StatusOK, map[string]any{
		"device_name": "Desk Rig",
	})
	devices := getJSON(t, router, "/api/devices", channelID, http.StatusOK)
	deviceItems := dataArray(t, devices)
	assertLen(t, deviceItems, 1, "devices")
	assertField(t, deviceItems[0], "name", "Desk Rig")

	syncHistory := getJSON(t, router, "/api/sync/history?limit=2", channelID, http.StatusOK)
	assertDataBool(t, syncHistory, "has_more", true)
	assertLen(t, dataItems(t, syncHistory), 2, "sync history first page")

	stats := getJSON(t, router, "/api/stats", channelID, http.StatusOK)
	assertNumberAtLeast(t, stats, "clipboard_item_count", 3)
	assertNumberAtLeast(t, stats, "total_devices", 1)
	assertNumberAtLeast(t, stats, "sync_count", 3)

	doJSON(t, router, http.MethodPost, "/api/clipboard", channelID, http.StatusBadRequest, map[string]any{
		"title":          "Bad actor",
		"content":        "Should fail before save",
		"type":           "text",
		"device_id":      "unregistered-device",
		"device_type":    "desktop",
		"content_format": "plain",
	})
	doJSON(t, router, http.MethodPost, "/api/sync/log", channelID, http.StatusOK, map[string]any{
		"device_id": deviceID,
		"content":   "manual sync log",
	})
	doJSON(t, router, http.MethodPost, "/api/sync/log", channelID, http.StatusBadRequest, map[string]any{
		"deviceId": deviceID,
		"content":  "old camelCase should fail",
	})

	deleteJSON(t, router, "/api/clipboard/"+stringField(t, codeItem, "id"), channelID, http.StatusOK, map[string]any{
		"device_id": deviceID,
	})
	statsAfterDelete := getJSON(t, router, "/api/stats", channelID, http.StatusOK)
	assertNumberAtLeast(t, statsAfterDelete, "clipboard_item_count", 2)
}

func newMainFlowRouter(t *testing.T) *gin.Engine {
	t.Helper()

	t.Setenv("HOME", t.TempDir())
	t.Setenv("USERPROFILE", t.TempDir())
	gin.SetMode(gin.TestMode)

	router, err := BuildRouterWithConfig(&config.Config{
		Host: "127.0.0.1",
		Port: 8080,
		Security: config.SecurityConfig{
			MaxBodyBytes: 2 << 20,
		},
		Log: config.LogConfig{SQL: "silent"},
	})
	if err != nil {
		t.Fatalf("build router: %v", err)
	}
	t.Cleanup(func() {
		sqlDB, err := db.GetDB().DB()
		if err != nil {
			t.Fatalf("get sql db: %v", err)
		}
		if err := sqlDB.Close(); err != nil {
			t.Fatalf("close db: %v", err)
		}
	})
	return router
}

func saveClipboardViaAPI(t *testing.T, router http.Handler, channelID, deviceID, content, itemType, title string) map[string]any {
	t.Helper()

	return postJSON(t, router, "/api/clipboard", channelID, http.StatusOK, map[string]any{
		"title":          title,
		"content":        content,
		"type":           itemType,
		"device_id":      deviceID,
		"device_type":    "desktop",
		"content_format": "plain",
	})
}

func getJSON(t *testing.T, router http.Handler, target, channelID string, expectedStatus int) testAPIResponse {
	t.Helper()
	return doJSON(t, router, http.MethodGet, target, channelID, expectedStatus, nil)
}

func postJSON(t *testing.T, router http.Handler, target, channelID string, expectedStatus int, body any) map[string]any {
	t.Helper()
	response := doJSON(t, router, http.MethodPost, target, channelID, expectedStatus, body)
	return dataObject(t, response)
}

func putJSON(t *testing.T, router http.Handler, target, channelID string, expectedStatus int, body any) map[string]any {
	t.Helper()
	response := doJSON(t, router, http.MethodPut, target, channelID, expectedStatus, body)
	return dataObject(t, response)
}

func deleteJSON(t *testing.T, router http.Handler, target, channelID string, expectedStatus int, body any) testAPIResponse {
	t.Helper()
	return doJSON(t, router, http.MethodDelete, target, channelID, expectedStatus, body)
}

func doJSON(t *testing.T, router http.Handler, method, target, channelID string, expectedStatus int, body any) testAPIResponse {
	t.Helper()

	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}

	request := httptest.NewRequest(method, target, &requestBody)
	request.Header.Set("Content-Type", "application/json")
	if channelID != "" {
		request.Header.Set("X-Channel-ID", channelID)
	}
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != expectedStatus {
		t.Fatalf("%s %s expected status %d, got %d body=%s", method, target, expectedStatus, response.Code, response.Body.String())
	}

	var apiResponse testAPIResponse
	if err := json.Unmarshal(response.Body.Bytes(), &apiResponse); err != nil {
		t.Fatalf("decode response for %s %s: %v body=%s", method, target, err, response.Body.String())
	}
	if expectedStatus == http.StatusOK && !apiResponse.Success {
		t.Fatalf("%s %s expected success response, got %#v", method, target, apiResponse)
	}
	return apiResponse
}

func dataObject(t *testing.T, response testAPIResponse) map[string]any {
	t.Helper()

	var data map[string]any
	if err := json.Unmarshal(response.Data, &data); err != nil {
		t.Fatalf("decode data object: %v raw=%s", err, string(response.Data))
	}
	return data
}

func dataArray(t *testing.T, response testAPIResponse) []map[string]any {
	t.Helper()

	var data []map[string]any
	if err := json.Unmarshal(response.Data, &data); err != nil {
		t.Fatalf("decode data array: %v raw=%s", err, string(response.Data))
	}
	return data
}

func dataItems(t *testing.T, response testAPIResponse) []map[string]any {
	t.Helper()

	data := dataObject(t, response)
	rawItems, ok := data["items"]
	if !ok {
		t.Fatalf("expected data.items in %#v", data)
	}
	items, ok := rawItems.([]any)
	if !ok {
		t.Fatalf("expected data.items array, got %T", rawItems)
	}

	result := make([]map[string]any, 0, len(items))
	for _, item := range items {
		itemObject, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("expected item object, got %T", item)
		}
		result = append(result, itemObject)
	}
	return result
}

func assertDataField(t *testing.T, data map[string]any, key, expected string) {
	t.Helper()
	assertField(t, data, key, expected)
}

func assertField(t *testing.T, data map[string]any, key, expected string) {
	t.Helper()

	actual, ok := data[key].(string)
	if !ok {
		t.Fatalf("expected %s to be string in %#v", key, data)
	}
	if actual != expected {
		t.Fatalf("expected %s=%q, got %q", key, expected, actual)
	}
}

func assertDataBool(t *testing.T, response testAPIResponse, key string, expected bool) {
	t.Helper()

	data := dataObject(t, response)
	actual, ok := data[key].(bool)
	if !ok {
		t.Fatalf("expected %s to be bool in %#v", key, data)
	}
	if actual != expected {
		t.Fatalf("expected %s=%v, got %v", key, expected, actual)
	}
}

func assertNumberAtLeast(t *testing.T, response testAPIResponse, key string, min float64) {
	t.Helper()

	data := dataObject(t, response)
	actual, ok := data[key].(float64)
	if !ok {
		t.Fatalf("expected %s to be number in %#v", key, data)
	}
	if actual < min {
		t.Fatalf("expected %s >= %.0f, got %.0f", key, min, actual)
	}
}

func stringField(t *testing.T, data map[string]any, key string) string {
	t.Helper()
	return stringFieldFromMap(t, data, key)
}

func stringFieldFromMap(t *testing.T, data map[string]any, key string) string {
	t.Helper()

	value, ok := data[key].(string)
	if !ok {
		t.Fatalf("expected %s to be string in %#v", key, data)
	}
	return value
}

func assertLen[T any](t *testing.T, values []T, expected int, label string) {
	t.Helper()

	if len(values) != expected {
		t.Fatalf("expected %s length %d, got %d: %#v", label, expected, len(values), values)
	}
}
