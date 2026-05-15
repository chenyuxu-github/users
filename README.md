# 用户管理系统

一个基于 Go、Gin、GORM 和 MySQL 的简单用户管理项目，提供用户的增删改查接口，并内置一个无需构建工具的前端页面。

## 功能

- 用户列表查询
- 新建用户
- 编辑用户
- 删除用户
- 左侧菜单布局
- 新建和编辑弹窗表单

## 技术栈

- Go
- Gin
- GORM
- Viper
- MySQL
- HTML / CSS / JavaScript

## 目录结构

```text
.
├── cmd
│   └── server
│       └── main.go          # 服务入口
├── config
│   └── config.toml          # 项目配置
├── internal
│   ├── config
│   │   └── config.go        # 配置读取
│   ├── database
│   │   └── db.go            # 数据库连接
│   ├── handler
│   │   └── user.go          # 用户接口处理
│   ├── model
│   │   └── user.go          # 用户模型
│   └── router
│       └── router.go        # 路由注册
├── web
│   ├── index.html           # 前端页面
│   └── assets
│       ├── app.js           # 前端交互逻辑
│       └── styles.css       # 页面样式
├── go.mod
├── go.sum
└── README.md
```

## 数据库配置

数据库连接和服务监听地址通过 Viper 从 `config/config.toml` 读取：

```toml
[server]
host = "0.0.0.0"
port = 8080

[database]
host = "localhost"
port = 3306
username = "root"
password = "123"
name = "test"
charset = "utf8mb4"
parse_time = true
loc = "Local"
```

也可以使用环境变量覆盖配置：

```bash
SERVER_PORT=9090 DATABASE_NAME=test go run ./cmd/server
```

请先确认本地 MySQL 中存在 `test` 数据库。项目启动时会通过 GORM 自动创建或迁移 `users` 表。

示例创建数据库 SQL：

```sql
CREATE DATABASE IF NOT EXISTS test DEFAULT CHARACTER SET utf8mb4;
```

## 启动项目

安装依赖：

```bash
go mod tidy
```

启动服务：

```bash
go run ./cmd/server
```

服务默认监听：

```text
http://localhost:8080
```

前端页面：

```text
http://localhost:8080/
```

## API 接口

接口统一前缀：

```text
/api/v1
```

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/api/v1/users` | 查询所有用户 |
| GET | `/api/v1/users/:id` | 查询单个用户 |
| POST | `/api/v1/users` | 创建用户 |
| PUT | `/api/v1/users/:id` | 更新用户 |
| DELETE | `/api/v1/users/:id` | 删除用户 |

创建或更新用户请求体示例：

```json
{
  "name": "张三",
  "email": "zhangsan@example.com"
}
```

## 常用命令

运行测试：

```bash
go test ./...
```

构建服务：

```bash
go build -o server ./cmd/server
```

启动构建后的服务：

```bash
./server
```
