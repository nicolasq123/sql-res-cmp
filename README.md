# sql-res-cmp

SQL 查询结果对比工具，支持 MySQL 和 ClickHouse。

## 功能

- 同时执行两个数据库查询
- 对比查询结果是否一致
- 支持按 key 列对比（用于无序结果的比较）

## 安装

```bash
go build -o diffq ./cmd
```

## 使用方法

```bash
diffq -d1 "<DSN1>" -q1 "<查询1>" -d2 "<DSN2>" -q2 "<查询2>" [-timeout 60s] [-key "col1,col2"]
```

### DSN 格式

**MySQL:**
```
mysql://user:pass@host:port/dbname?timeout=30s
```

**ClickHouse:**
```
clickhouse://host:port/dbname?username=user&password=pass&timeout=30s
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
| `-timeout` | 查询超时时间 | 60s |
| `-key` | 用于对比的 key 列（逗号分隔） | 无 |
| `-ding` | 钉钉机器人 webhook 地址，比对不一致时自动发送告警 | 无 |

### 返回值

- `0`: 结果一致
- `1`: 结果不一致

## 项目结构

```
.
├── cmd/            # CLI 入口
├── comparator/     # 结果比较逻辑
├── executor/       # 数据库执行器
│   ├── executor.go      # 接口定义和 DSN 解析
│   ├── mysql.go         # MySQL 实现
│   └── clickhouse.go    # ClickHouse 实现
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
