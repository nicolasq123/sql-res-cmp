package alarm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const baseUrl = "https://oapi.dingtalk.com/robot/send?access_token="

type DingAlarm struct {
	webhook string
	client  *http.Client
}

func NewDingAlarm(token string) *DingAlarm {
	if token == "" {
		panic("token is required")
	}

	return &DingAlarm{webhook: baseUrl + token, client: &http.Client{Timeout: 10 * time.Second}}
}

// SendAlert implements AlertSender interface
func (d *DingAlarm) Send(ctx context.Context, name, q1, q2, diffText string) error {
	type DingMsg struct {
		MsgType  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
	}

	var sb strings.Builder
	sb.WriteString(alarmPrefix)
	sb.WriteString("### ⚠️ SQL Result Mismatch\n\n")

	if name != "" {
		sb.WriteString(fmt.Sprintf("**Name:** %s\n\n", name))
	}

	sb.WriteString("---\n\n")
	sb.WriteString("#### 📊 Summary\n")
	sb.WriteString(fmt.Sprintf("- Time: `%s`\n", time.Now().Format("2006-01-02 15:04:05")))

	// Parse and extract summary from diffText
	if strings.Contains(diffText, "Results mismatch") {
		parts := strings.SplitN(diffText, "\n", 2)
		if len(parts) > 0 {
			sb.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(parts[0])))
		}
	}

	sb.WriteString("\n---\n\n")

	if q1 != "" {
		sb.WriteString("#### 🔍 Query 1\n")
		sb.WriteString("```sql\n")
		sb.WriteString(truncate(q1, 500))
		sb.WriteString("\n```\n\n")
	}

	if q2 != "" {
		sb.WriteString("#### 🔍 Query 2\n")
		sb.WriteString("```sql\n")
		sb.WriteString(truncate(q2, 500))
		sb.WriteString("\n```\n\n")
	}

	sb.WriteString("---\n\n")
	sb.WriteString("#### 📝 Diff Details\n")
	sb.WriteString("```\n")
	sb.WriteString(truncate(diffText, 2000))
	sb.WriteString("\n```\n")

	var msg DingMsg
	msg.MsgType = "markdown"
	msg.Markdown.Title = "⚠️ SQL Result Mismatch"
	msg.Markdown.Text = sb.String()

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("do: %w", err)
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status: %d", resp.StatusCode)
	}
	return nil
}

// func truncate(s string, max int) string {
// 	if len(s) <= max {
// 		return s
// 	}
// 	return s[:max] + "..."
// }
