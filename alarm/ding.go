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

	"github.com/nicolasq123/sql-res-cmp/comparator"
)

// DingAlarm 钉钉告警器
type DingAlarm struct {
	webhook string
	client  *http.Client
}

// NewDingAlarm 创建钉钉告警器
func NewDingAlarm(webhook string) *DingAlarm {
	return &DingAlarm{
		webhook: webhook,
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

// Send 发送差异告警
func (d *DingAlarm) Send(ctx context.Context, diff *comparator.Diff, d1, q1, d2, q2 string) error {
	// 钉钉markdown消息格式
	type DingMsg struct {
		MsgType string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
	}

	var msg DingMsg
	msg.MsgType = "markdown"
	msg.Markdown.Title = "SQL 结果比对不一致告警"

	// 构建消息内容
	var sb strings.Builder
	sb.WriteString("### SQL 结果比对不一致\n\n")
	sb.WriteString("#### 基本信息\n")
	sb.WriteString(fmt.Sprintf("- 数据库1: `%s`\n", d1))
	sb.WriteString(fmt.Sprintf("- 查询1: `%s`\n", truncateString(q1, 100)))
	sb.WriteString(fmt.Sprintf("- 数据库2: `%s`\n", d2))
	sb.WriteString(fmt.Sprintf("- 查询2: `%s`\n", truncateString(q2, 100)))
	sb.WriteString(fmt.Sprintf("- 比对时间: `%s`\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("#### 差异详情\n")
	if len(diff.LeftOnly) > 0 {
		sb.WriteString(fmt.Sprintf("- 仅左侧有 %d 行\n", len(diff.LeftOnly)))
	}
	if len(diff.RightOnly) > 0 {
		sb.WriteString(fmt.Sprintf("- 仅右侧有 %d 行\n", len(diff.RightOnly)))
	}
	if len(diff.Modified) > 0 {
		sb.WriteString(fmt.Sprintf("- 不一致 %d 行\n", len(diff.Modified)))
	}

	// 列名
	sb.WriteString("\n#### 列名\n")
	sb.WriteString("```\n")
	sb.WriteString(strings.Join(diff.Columns, ", "))
	sb.WriteString("\n```\n")

	// 最多展示3行差异
	if len(diff.Modified) > 0 {
		sb.WriteString("\n#### 不一致示例\n")
		for i, m := range diff.Modified {
			if i >= 3 {
				sb.WriteString(fmt.Sprintf("\n... 还有 %d 行差异未显示", len(diff.Modified)-3))
				break
			}
			sb.WriteString(fmt.Sprintf("第 %d 行:\n", i+1))
			sb.WriteString("- 左侧: `")
			sb.WriteString(strings.Join(m.Left, ", "))
			sb.WriteString("`\n")
			sb.WriteString("- 右侧: `")
			sb.WriteString(strings.Join(m.Right, ", "))
			sb.WriteString("`\n")
		}
	}

	msg.Markdown.Text = sb.String()

	// 发送请求
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", d.webhook, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	// 完全读取响应体以确保连接可以复用
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dingding api returned: %d", resp.StatusCode)
	}

	return nil
}

// truncateString 字符串截断
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
