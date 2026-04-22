package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nicolasq123/sql-res-cmp/alarm"
	"github.com/nicolasq123/sql-res-cmp/comparator"
	"github.com/nicolasq123/sql-res-cmp/executor"
)

func main() {
	var d1, q1, d2, q2 string
	var timeout time.Duration
	var keyCols string
	var dingWebhook string

	flag.StringVar(&d1, "d1", "", "第一个数据库 DSN")
	flag.StringVar(&q1, "q1", "", "第一个查询")
	flag.StringVar(&d2, "d2", "", "第二个数据库 DSN")
	flag.StringVar(&q2, "q2", "", "第二个查询")
	flag.DurationVar(&timeout, "timeout", 60*time.Second, "查询超时")
	flag.StringVar(&keyCols, "key", "", "用于比较的 key 列，逗号分隔")
	flag.StringVar(&dingWebhook, "ding", "", "钉钉机器人 webhook 地址，比对不一致时发送告警")

	flag.Parse()

	if d1 == "" || q1 == "" || d2 == "" || q2 == "" {
		fmt.Println("参数错误: 需要指定 -d1 -q1 -d2 -q2")
		flag.Usage()
		os.Exit(1)
	}

	db1, err := executor.NewDB(d1)
	if err != nil {
		fmt.Printf("连接数据库1错误: %v\n", err)
		os.Exit(1)
	}
	defer db1.Close()

	db2, err := executor.NewDB(d2)
	if err != nil {
		fmt.Printf("连接数据库2错误: %v\n", err)
		os.Exit(1)
	}
	defer db2.Close()

	cols1, rows1, err := query(db1, q1)
	if err != nil {
		fmt.Printf("执行查询1错误: %v\n", err)
		os.Exit(1)
	}

	_, rows2, err := query(db2, q2)
	if err != nil {
		fmt.Printf("执行查询2错误: %v\n", err)
		os.Exit(1)
	}

	c := comparator.NewComparator()
	var diff *comparator.Diff

	if keyCols != "" {
		keys := strings.Split(keyCols, ",")
		diff = c.CompareByKey(rows1, rows2, keys, cols1)
	} else {
		diff = c.Compare(rows1, rows2, cols1)
	}

	fmt.Print(diff.String())
	if !diff.IsEmpty() {
		if dingWebhook != "" {
			dingAlarm := alarm.NewDingAlarm(dingWebhook)
			alarmCtx, alarmCancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer alarmCancel()
			if err := dingAlarm.Send(alarmCtx, diff, d1, q1, d2, q2); err != nil {
				fmt.Printf("发送钉钉告警失败: %v\n", err)
			} else {
				fmt.Println("钉钉告警发送成功")
			}
		}
		os.Exit(1)
	}
}

func query(db *sqlx.DB, q string) ([]string, [][]string, error) {
	rows, err := db.Queryx(q)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var result [][]string
	for rows.Next() {
		dests := make([]any, len(cols))
		byteSlices := make([][]byte, len(cols))
		for i := range dests {
			dests[i] = &byteSlices[i]
		}
		if err := rows.Scan(dests...); err != nil {
			return nil, nil, err
		}
		row := make([]string, len(cols))
		for i, b := range byteSlices {
			if b == nil {
				row[i] = "NULL"
			} else {
				row[i] = string(b)
			}
		}
		result = append(result, row)
	}
	return cols, result, rows.Err()
}
