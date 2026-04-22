# sql-res-cmp

SQL 查询结果对比工具，支持 MySQL/ClickHouse/PostgreSQL。

[English Documentation](README.md)

## 功能

- 同时执行两个数据库查询
- 对比查询结果是否一致
- 支持按 key 列对比（用于无序结果的比较）
- 比对不一致时支持钉钉/WeChat 告警
  - DingTalk 通知【使用 webhook】
  - WeChat 通知【调用 weclaw 生成的 credential】

## 安装

```bash
go build -o diffq ./cmd
```

## 使用方法

```bash
diffq -d1 "<DSN1>" -q1 "<query1>" -d2 "<DSN2>" -q2 "<query2>" [-key "col1,col2"] [-ding "webhook"] [-wechat true] [-name myscript]
```

### DSN 格式

**MySQL:**
```
mysql://user:pass@host:port/dbname
```

**ClickHouse:**
```
clickhouse://host:port/dbname?username=user&password=pass
```

**PostgreSQL:**
```
postgres://user:pass@host:port/dbname
```

### 示例

**相同查询对比:**
```bash
./diffq \
  -d1 "mysql://root:123456@127.0.0.1:3306/testdb" \
  -q1 "SELECT id, name FROM users WHERE status=1" \
  -d2 "mysql://root:123456@127.0.0.1:3306/testdb" \
  -q2 "SELECT id, name FROM users WHERE status=1"
```

**比对失败发送钉钉告警:**
```bash
./diffq \
  -d1 "mysql://root:123456@127.0.0.1:3306/testdb" \
  -q1 "SELECT id, name FROM users WHERE status=1" \
  -d2 "mysql://root:123456@127.0.0.1:3306/testdb" \
  -q2 "SELECT id, name FROM users WHERE status=1" \
  -ding "https://oapi.dingtalk.com/robot/send?access_token=xxx"
```

**跨数据库对比:**
```bash
./diffq \
  -d1 "mysql://root:123456@127.0.0.1:3306/mysql" \
  -q1 "SELECT 1 as id, 'hello' as name" \
  -d2 "clickhouse://localhost:9000/default?username=default&password=" \
  -q2 "SELECT 1 as id, 'hello' as name"
```

**使用 key 列对比（结果无序时）:**
```bash
./diffq \
  -d1 "mysql://..." -q1 "SELECT * FROM orders" \
  -d2 "mysql://..." -q2 "SELECT * FROM orders ORDER BY id" \
  -key "id"
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-d1` | 第一个数据库 DSN | 必需 |
| `-q1` | 第一个查询 SQL | 必需 |
| `-d2` | 第二个数据库 DSN | 必需 |
| `-q2` | 第二个查询 SQL | 必需 |
| `-key` | 用于对比的 key 列（逗号分隔） | 无 |
| `-ding` | 钉钉机器人 webhook 地址 | 无 |
| `-wechat` | 发送 WeChat 告警 | false |
| `-timeout` | 超时时间 | 120s |
| `-name` | 对比任务名称 | 无 |

### 返回值

- `0`: 结果一致
- `1`: 结果不一致

## 项目结构

```
.
├── cmd/            # CLI 入口
├── cmp/            # 比较和执行逻辑
├── alarm/          # 告警（钉钉、WeChat）
├── .github/        # CI 配置
├── go.mod
└── README.md
```

## 测试

```bash
go test ./... -v
go test ./... -cover
```

## License

MIT
