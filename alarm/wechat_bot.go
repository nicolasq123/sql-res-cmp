package alarm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type WeChatBot struct {
	creds   map[string]any
	baseURL string
}

// 使用weclaw生成的credential
// 参考: wechatbot
func NewWeChatBot(credentials map[string]any) (*WeChatBot, error) {
	if credentials == nil {
		var err error
		credentials, err = loadCredentials()
		if err != nil {
			return nil, err
		}
	}

	baseURL := "https://ilinkai.weixin.qq.com"
	if u, ok := credentials["baseurl"].(string); ok && u != "" {
		baseURL = u
	}

	return &WeChatBot{
		creds:   credentials,
		baseURL: baseURL,
	}, nil
}

func loadCredentials() (map[string]any, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get user home: %w", err)
	}

	accountsDir := filepath.Join(homeDir, ".weclaw", "accounts")
	entries, err := os.ReadDir(accountsDir)
	if err != nil {
		return nil, fmt.Errorf("read accounts dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) != ".json" || len(name) > 9 && name[len(name)-9:] == ".sync.json" {
			continue
		}

		path := filepath.Join(accountsDir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		var creds map[string]any
		if err := json.Unmarshal(data, &creds); err != nil {
			continue
		}

		if token, ok := creds["bot_token"].(string); ok && token != "" {
			return creds, nil
		}
	}

	return nil, fmt.Errorf("no credentials found, please login first")
}

func generateWechatUin() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rng.Uint32()
	s := fmt.Sprintf("%d", n)
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func (w *WeChatBot) SendText(ctx context.Context, text string) error {
	return w.Send(ctx, "", "", "", text)
}

func (w *WeChatBot) Send(ctx context.Context, name, q1, q2, diffText string) error {
	url := w.baseURL + "/ilink/bot/sendmessage"

	var sb strings.Builder
	sb.WriteString("⚠️ SQL Result Mismatch\n")
	sb.WriteString("═══════════════════════\n\n")

	if name != "" {
		sb.WriteString(fmt.Sprintf("Name: %s\n\n", name))
	}

	sb.WriteString("📊 Summary\n")
	sb.WriteString(fmt.Sprintf("Time: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	if strings.Contains(diffText, "Results mismatch") {
		parts := strings.SplitN(diffText, "\n", 2)
		if len(parts) > 0 {
			sb.WriteString(fmt.Sprintf("%s\n", strings.TrimSpace(parts[0])))
		}
	}

	sb.WriteString("\n═══════════════════════\n\n")

	if q1 != "" {
		sb.WriteString("🔍 Query 1\n")
		sb.WriteString(truncate(q1, 300))
		sb.WriteString("\n\n")
	}

	if q2 != "" {
		sb.WriteString("🔍 Query 2\n")
		sb.WriteString(truncate(q2, 300))
		sb.WriteString("\n\n")
	}

	sb.WriteString("═══════════════════════\n\n")
	sb.WriteString("📝 Diff Details\n")
	sb.WriteString(truncate(diffText, 1500))

	text := alarmPrefix + sb.String()

	msg := map[string]any{
		"from_user_id":  w.creds["ilink_bot_id"],
		"to_user_id":    w.creds["ilink_user_id"],
		"client_id":     uuid.New().String(),
		"message_type":  2,
		"message_state": 2,
		"item_list": []any{
			map[string]any{
				"type": 1,
				"text_item": map[string]any{
					"text": text,
				},
			},
		},
		"context_token": "",
	}

	data := map[string]any{
		"msg":       msg,
		"base_info": map[string]any{},
	}

	body, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AuthorizationType", "ilink_bot_token")
	req.Header.Set("Authorization", "Bearer "+w.creds["bot_token"].(string))
	req.Header.Set("X-WECHAT-UIN", generateWechatUin())

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if ret, ok := result["ret"].(float64); ok && ret != 0 {
			errmsg, _ := result["errmsg"].(string)
			return fmt.Errorf("send failed: ret=%v errmsg=%s", ret, errmsg)
		}
	}

	return nil
}
