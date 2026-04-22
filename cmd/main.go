package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nicolasq123/sql-res-cmp/alarm"
	"github.com/nicolasq123/sql-res-cmp/cmp"
)

func main() {
	var name, d1, q1, d2, q2, keyCols, dingToken string
	var wechat bool
	var timeout time.Duration

	flag.StringVar(&name, "name", "", "Name")
	flag.StringVar(&d1, "d1", "", "Database 1 DSN")
	flag.StringVar(&q1, "q1", "", "Query 1")
	flag.StringVar(&d2, "d2", "", "Database 2 DSN")
	flag.StringVar(&q2, "q2", "", "Query 2")
	flag.StringVar(&keyCols, "key", "", "Key columns for comparison")
	flag.StringVar(&dingToken, "ding", "", "DingTalk token")
	flag.BoolVar(&wechat, "wechat", false, "Send WeChat alert")
	flag.DurationVar(&timeout, "timeout", 120*time.Second, "Timeout duration")
	flag.Parse()

	if d1 == "" || q1 == "" || d2 == "" || q2 == "" {
		fmt.Println("Required: -d1 -q1 -d2 -q2")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	db1, err := cmp.NewDB(d1)
	myPanic(err)
	defer db1.Close()

	db2, err := cmp.NewDB(d2)
	myPanic(err)
	defer db2.Close()

	cols1, rows1, err := cmp.Query(ctx, db1, q1)
	myPanic(err)

	_, rows2, err := cmp.Query(ctx, db2, q2)
	myPanic(err)

	var diff *cmp.Diff
	c := cmp.NewComparator()
	if keyCols != "" {
		diff = c.CompareByKey(rows1, rows2, strings.Split(keyCols, ","), cols1)
	} else {
		diff = c.Compare(rows1, rows2, cols1)
	}

	fmt.Print(diff.String())
	if !diff.IsEmpty() {
		alarmCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if dingToken != "" {
			if err := alarm.NewDingAlarm(dingToken).Send(alarmCtx, name, q1, q2, diff.String()); err != nil {
				fmt.Printf("ding: %v\n", err)
			}
		}
		if wechat {
			bot, err := alarm.NewWeChatBot(nil)
			if err != nil {
				fmt.Printf("wechat bot: %v\n", err)
			} else if err := bot.Send(alarmCtx, name, q1, q2, diff.String()); err != nil {
				fmt.Printf("wechat: %v\n", err)
			}
		}
	}
}

func myPanic(err error) {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		panic(err)
	}
}
