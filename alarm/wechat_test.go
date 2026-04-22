package alarm

import (
	"context"
	"testing"
	"time"
)

func TestWeChatBotSendText(t *testing.T) {
	bot, err := NewWeChatBot(nil)
	if err != nil {
		t.Fatalf("create bot failed: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	err = bot.SendText(ctx, "Hello from Go test!")
	if err != nil {
		t.Fatalf("send text failed: %v", err)
	}
	t.Log("send success")
}
