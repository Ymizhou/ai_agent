package user_test

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"aicode/cmd"
	"aicode/internal/model/entity"
)

func init() {
	// memstore 序列化 entity.User 需要注册 gob 类型
	gob.Register(&entity.User{})

	// providers.go 中 LoadConfig 使用相对路径 "config.yml"，
	// 测试从 test/user/ 执行，切换到项目根目录使路径正确
	if err := os.Chdir("../../"); err != nil {
		panic(fmt.Sprintf("切换工作目录失败: %v", err))
	}
}

// apiResp 通用响应结构
type apiResp[T any] struct {
	Code    int    `json:"code"`
	Data    T      `json:"data"`
	Message string `json:"message"`
}

// newTestServer 初始化完整应用并启动 httptest.Server
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	app, err := cmd.InitializeApp()
	if err != nil {
		t.Fatalf("初始化应用失败: %v", err)
	}
	return httptest.NewServer(app.Router)
}

// doJSON 发送 JSON 请求，可携带可选 cookie
func doJSON(t *testing.T, client *http.Client, method, url string, body any, cookies ...*http.Cookie) *http.Response {
	t.Helper()
	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("序列化请求体失败: %v", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		t.Fatalf("创建请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	for _, c := range cookies {
		if c != nil {
			req.AddCookie(c)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("发送请求失败: %v", err)
	}
	return resp
}

// decodeResp 将响应体解码到 apiResp[T]
func decodeResp[T any](t *testing.T, resp *http.Response) apiResp[T] {
	t.Helper()
	defer resp.Body.Close()
	var result apiResp[T]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
	return result
}

// sessionCookie 从响应中提取名为 session_id 的 cookie
func sessionCookie(resp *http.Response) *http.Cookie {
	for _, c := range resp.Cookies() {
		if c.Name == "session_id" {
			return c
		}
	}
	return nil
}

// randAccount 生成随机账号，避免重复注册冲突
func randAccount() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("testuser%06d", r.Int31n(1000000))
}

// TestUserFlow 覆盖：注册 → 登录 → 获取当前用户 → 登出 → 登出后访问被拦截
func TestUserFlow(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	client := ts.Client()
	base := ts.URL + "/api/v1"

	account := randAccount()
	password := "password123"

	var sessionCook *http.Cookie

	// ── 1. 注册 ──────────────────────────────────────────────────────────
	t.Run("register", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodPost, base+"/user/register", map[string]string{
			"userAccount":   account,
			"userPassword":  password,
			"checkPassword": password,
		})
		result := decodeResp[int64](t, resp)

		if result.Code != 0 {
			t.Fatalf("注册失败，code=%d msg=%s", result.Code, result.Message)
		}
		if result.Data <= 0 {
			t.Fatalf("注册成功但返回的用户 ID 无效: %d", result.Data)
		}
		t.Logf("注册成功，用户 ID: %d", result.Data)
	})

	// ── 2. 用错误密码登录（预期失败）────────────────────────────────────
	t.Run("login_wrong_password", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodPost, base+"/user/login", map[string]string{
			"userAccount":  account,
			"userPassword": "wrongpassword",
		})
		result := decodeResp[any](t, resp)

		if result.Code == 0 {
			t.Fatal("错误密码登录应失败，但返回了 code=0")
		}
		t.Logf("错误密码登录被正确拒绝，code=%d msg=%s", result.Code, result.Message)
	})

	// ── 3. 正确登录 ──────────────────────────────────────────────────────
	t.Run("login", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodPost, base+"/user/login", map[string]string{
			"userAccount":  account,
			"userPassword": password,
		})

		result := decodeResp[map[string]any](t, resp)
		if result.Code != 0 {
			t.Fatalf("登录失败，code=%d msg=%s", result.Code, result.Message)
		}

		sessionCook = sessionCookie(resp)
		if sessionCook == nil {
			t.Fatal("登录成功但响应中未包含 session_id cookie")
		}
		t.Logf("登录成功，session cookie 已获取: %s=%.12s...", sessionCook.Name, sessionCook.Value)
	})

	t.Run("login2", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodPost, base+"/user/login", map[string]string{
			"userAccount":  "admin",
			"userPassword": "12345678",
		})

		result := decodeResp[map[string]any](t, resp)
		if result.Code != 0 {
			t.Fatalf("登录失败，code=%d msg=%s", result.Code, result.Message)
		}

		sessionCook = sessionCookie(resp)
		if sessionCook == nil {
			t.Fatal("登录成功但响应中未包含 session_id cookie")
		}
		t.Logf("登录成功，session cookie 已获取: %s=%.12s...", sessionCook.Name, sessionCook.Value)
	})

	// ── 4. 携带 session cookie 获取当前登录用户 ──────────────────────────
	t.Run("get_login_user", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodGet, base+"/user/get/login", nil, sessionCook)
		result := decodeResp[map[string]any](t, resp)

		if result.Code != 0 {
			t.Fatalf("获取登录用户失败，code=%d msg=%s", result.Code, result.Message)
		}
		t.Logf("当前登录用户: userAccount=%v", result.Data["userAccount"])
	})

	// ── 5. 不带 cookie 访问受保护接口（预期被拦截）────────────────────────
	t.Run("access_without_cookie", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodGet, base+"/user/get/login", nil)
		result := decodeResp[any](t, resp)

		if result.Code == 0 {
			t.Fatal("未携带 cookie 访问受保护接口应返回未登录错误，但返回了 code=0")
		}
		t.Logf("未登录访问被正确拦截，code=%d msg=%s", result.Code, result.Message)
	})

	// ── 6. 登出 ──────────────────────────────────────────────────────────
	t.Run("logout", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodPost, base+"/user/logout", nil, sessionCook)
		result := decodeResp[bool](t, resp)

		if result.Code != 0 {
			t.Fatalf("登出失败，code=%d msg=%s", result.Code, result.Message)
		}
		if !result.Data {
			t.Fatal("登出返回 code=0 但 data=false")
		}
		t.Logf("登出成功")
	})

	// ── 7. 登出后再访问受保护接口（预期被拦截）────────────────────────────
	t.Run("access_after_logout", func(t *testing.T) {
		resp := doJSON(t, client, http.MethodGet, base+"/user/get/login", nil, sessionCook)
		result := decodeResp[any](t, resp)

		if result.Code == 0 {
			t.Fatal("登出后访问受保护接口应返回未登录错误，但返回了 code=0")
		}
		t.Logf("登出后访问被正确拦截，code=%d msg=%s", result.Code, result.Message)
	})
}
