package alarm

import "context"

var alarmPrefix = "[diffq_alarm]: "

// AlertSender is the interface for sending alerts
type AlertSender interface {
	SendAlert(ctx context.Context, name, q1, q2, diffText string) error
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

// func SetAlarmPrefix(prefix string) {
// 	alarmPrefix = prefix
// }

// func GetAlarmPrefix() string {
// 	return alarmPrefix
// }
