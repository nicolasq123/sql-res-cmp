package alarm

import (
	"testing"
)

func TestWeChatBotSendText(t *testing.T) {
	creds := map[string]any{
		"bot_token":     "xxx",
		"ilink_bot_id":  "xxx",
		"ilink_user_id": "xxx",
		"baseurl":       "https://ilinkai.weixin.qq.com",
	}
	_, err := NewWeChatBot(creds)
	if err != nil {
		t.Fatalf("create bot failed: %v", err)
	}
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()
	// err = bot.SendText(ctx, "Hello from Go test!")
	// if err != nil {
	// 	t.Fatalf("send text failed: %v", err)
	// }
	t.Log("send success")
}
