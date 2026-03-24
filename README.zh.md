# Token Manager - LLM Token 管理系统

一个用于管理员工信息和 LLM Token 分配的 REST API 后端服务。该系统允许管理员管理员工、发放带有配额限制的 Token、跟踪 Token 使用情况并生成统计数据。

## 目录

1. [概述](#概述)
2. [功能特性](#功能特性)
3. [架构设计](#架构设计)
4. [项目结构](#项目结构)
5. [环境要求](#环境要求)
6. [安装部署](#安装部署)
7. [运行应用程序](#运行应用程序)
8. [API 文档](#api-文档)
9. [前端配置](#前端配置)
10. [数据存储](#数据存储)
11. [API 示例](#api-示例)
12. [构建与测试](#构建与测试)
13. [代码风格指南](#代码风格指南)
14. [许可证](#许可证)

---

## 概述

Token Manager 是一个基于 Go 语言的 REST API 系统，用于在组织内管理 LLM（大型语言模型）Token。它提供了完整的解决方案：

- 员工注册与管理
- 带配额限制的 Token 发放
- Token 使用跟踪
- 使用统计与分析
- Token-员工关系映射

### 技术栈

- **后端**：Go 1.25+
- **Web 框架**：Gin
- **数据存储**：JSON 文件（嵌入式）
- **前端**：原生 HTML/CSS/JavaScript

---

## 功能特性

### 员工管理
- 创建、读取、更新和删除员工记录
- 跟踪员工部门与职位
- 存储员工联系方式

### Token 管理
- 发放带有可配置配额限制的 Token
- 设置 Token 过期时间（按天计算）
- 激活与撤销 Token
- 跟踪 Token 使用历史

### 统计与分析
- 系统级使用统计
- 员工级使用报告
- Token 利用率百分比
- 活跃与撤销 Token 数量统计

### Token-员工映射
- 查看所有 Token-员工关系
- 分页映射视图
- 剩余配额跟踪

---

## 架构设计

### 设计模式

该应用遵循多种设计模式：

1. **仓储模式**：`Store` 类型抽象了所有数据访问
2. **工厂模式**：`NewStore()` 作为工厂函数用于存储初始化
3. **处理器模式**：HTTP 处理器遵循统一的请求处理模式
4. **DTO 模式**：请求/响应 DTO 将 API 契约与内部模型分离

### 线程安全

所有数据访问都受到 `sync.RWMutex` 保护：
- 读操作使用 `RLock` 支持并发读取
- 写操作使用 `Lock` 进行独占访问
- 所有锁获取都使用 `defer` 确保自动解锁

### 数据流

```
客户端请求
    |
    v
Gin 路由器 --> 处理器 --> 存储方法 --> 数据存储 --> JSON 文件
                    |
                    v
              响应 DTO <-- 内部模型
```

---

## 项目结构

```
token-manager/
|
|-- [main.go](./main.go)          # 应用程序入口点
|-- [handlers.go](./handlers.go)      # HTTP 请求处理器
|-- [models.go](./models.go)        # 数据结构与 DTO
|-- [store.go](./store.go)         # 数据持久化与业务逻辑
|-- [data.json](./data.json)        # 数据存储文件（自动创建）
|-- [go.mod](./go.mod)           # Go 模块定义
|-- [go.sum](./go.sum)           # Go 模块校验和
|
|-- [index.html](./index.html)       # 前端入口点
|-- css/
|   |-- [styles.css](./css/styles.css)   # 所有 CSS 样式
|
|-- js/
|   |-- [app.js](./js/app.js)       # 前端 JavaScript
|
|-- [README.md](./README.md)        # 英文文档
|-- [README.zh.md](./README.zh.md)        # 本文件
```

---

## 环境要求

- **Go**：1.25 或更高版本
- **Web 浏览器**：支持 JavaScript 的现代浏览器（用于前端）
- **网络访问**：8080 端口可用

---

## 安装部署

### 1. 克隆或下载

```bash
# 进入项目目录
cd token-manager
```

### 2. 安装依赖

```bash
# 下载 Go 模块依赖
go mod tidy
```

这将自动下载 Gin Web 框架及其他所需包。

### 3. 验证安装

```bash
# 构建应用程序
go build -o token-manager.exe
```

---

## 运行应用程序

### 启动服务器

```bash
# 运行应用程序
go run main.go
```

服务器将在 `http://localhost:8080` 启动

### 预期输出

```
Server starting on :8080
```

### 数据文件初始化

首次运行时，应用程序将自动创建包含空数据结构的 `data.json` 文件：

```json
{
  "employees": [],
  "tokens": [],
  "usage_records": []
}
```

---

## API 文档

### 基础 URL

```
http://localhost:8080/api
```

### 员工端点

| 方法 | 端点 | 描述 |
|--------|----------|-------------|
| POST | /api/employees | 创建新员工 |
| GET | /api/employees | 获取所有员工 |
| GET | /api/employees/:id | 获取指定员工 |
| PUT | /api/employees/:id | 更新员工 |
| DELETE | /api/employees/:id | 删除员工 |

### Token 端点

| 方法 | 端点 | 描述 |
|--------|----------|-------------|
| POST | /api/tokens/issue | 发放新 Token |
| GET | /api/tokens | 获取所有 Token |
| GET | /api/employees/:id/tokens | 获取员工的 Token |
| POST | /api/tokens/:id/revoke | 撤销 Token |
| POST | /api/tokens/use | 记录 Token 使用 |
| GET | /api/tokens/:id/usage | 获取 Token 使用历史 |

### 统计端点

| 方法 | 端点 | 描述 |
|--------|----------|-------------|
| GET | /api/stats | 获取系统统计 |
| GET | /api/employees/:id/stats | 获取员工统计 |

### 映射端点

| 方法 | 端点 | 描述 |
|--------|----------|-------------|
| GET | /api/mappings | 获取 Token-员工映射 |

### /api/mappings 查询参数

| 参数 | 类型 | 默认值 | 描述 |
|-----------|------|---------|-------------|
| page | 整数 | 1 | 页码 |
| page_size | 整数 | 10 | 每页条目数 |

---

## 前端配置

### 运行前端

1. 确保后端服务器正在运行（`go run main.go`）
2. 在 Web 浏览器中打开 `index.html`
3. 前端将连接到 `http://localhost:8080/api`

### 前端功能

- **概览选项卡**：带系统统计的仪表板
- **员工选项卡**：创建和管理员工
- **Token 选项卡**：发放和管理 Token
- **映射选项卡**：查看 Token-员工关系

### 前端文件结构

```
[index.html](./index.html)       - 主 HTML 结构
[css/styles.css](./css/styles.css)   - 所有样式（深色主题、响应式）
[js/app.js](./js/app.js)        - API 通信与 UI 逻辑
```

---

## 数据存储

### 持久化

数据以 JSON 格式存储在 `data.json` 文件中。文件将自动：
- 首次运行时创建
- 后续运行时加载
- 每次写操作后更新

### 数据模式

```json
{
  "employees": [
    {
      "id": "字符串",
      "name": "字符串",
      "department": "字符串",
      "email": "字符串",
      "position": "字符串",
      "created_at": "时间戳",
      "updated_at": "时间戳"
    }
  ],
  "tokens": [
    {
      "id": "字符串",
      "employee_id": "字符串",
      "token_value": "字符串",
      "total_quota": "整数",
      "used_quota": "整数",
      "is_active": "布尔值",
      "issued_at": "时间戳",
      "expired_at": "时间戳",
      "revoked_at": "时间戳 | 空"
    }
  ],
  "usage_records": [
    {
      "id": "字符串",
      "token_id": "字符串",
      "used_at": "时间戳",
      "amount": "整数",
      "model": "字符串",
      "description": "字符串"
    }
  ]
}
```

---

## API 示例

### 创建员工

```bash
curl -X POST http://localhost:8080/api/employees \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@company.com",
    "department": "工程部",
    "position": "软件工程师"
  }'
```

### 发放 Token

```bash
curl -X POST http://localhost:8080/api/tokens/issue \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "1234567890",
    "total_quota": 1000000,
    "days_valid": 30
  }'
```

### 使用 Token

```bash
curl -X POST http://localhost:8080/api/tokens/use \
  -H "Content-Type: application/json" \
  -d '{
    "token_id": "token_123",
    "amount": 1500,
    "model": "gpt-4",
    "description": "代码审查"
  }'
```

### 获取系统统计

```bash
curl http://localhost:8080/api/stats
```

### 撤销 Token

```bash
curl -X POST http://localhost:8080/api/tokens/token_123/revoke
```

---

## 构建与测试

### 构建

```bash
# 构建应用程序
go build -o token-manager.exe

# 为不同平台构建
GOOS=windows GOARCH=amd64 go build -o token-manager.exe
GOOS=linux GOARCH=amd64 go build -o token-manager
```

### 运行

```bash
# 运行构建的可执行文件
./token-manager.exe

# 或直接使用 Go 运行
go run main.go
```

### 代码质量

```bash
# 格式化代码
go fmt ./...

# 运行 go vet
go vet ./...

# 运行所有测试（如果有）
go test ./...

# 运行单个测试
go test -run TestFunctionName ./...

# 带详细输出运行测试
go test -v ./...
```

### 依赖项

该项目使用单个外部依赖：

- **github.com/gin-gonic/gin** (v1.12.0) - Web 框架

在 [`go.mod`](./go.mod) 中查看所有依赖项。

---

## 代码风格指南

### Go 代码风格

本项目遵循标准 Go 约定：

1. **导入组织**：标准库优先，然后是外部包
2. **命名**：导出的标识符使用 PascalCase，私有标识符使用 camelCase
3. **错误处理**：返回错误而非 panic
4. **文档**：为所有导出的函数添加详细注释
5. **格式化**：使用 `go fmt` 自动格式化

### 前端代码风格

1. **CSS**：使用 CSS 自定义属性进行主题化
2. **JavaScript**：现代 ES6+ 特性
3. **分离**：CSS 和 JS 在单独的文件中
4. **注释**：JSDoc 风格的函数注释

### 文件组织

每个 Go 文件都有特定用途：

- [`main.go`](./main.go) - 应用程序入口点
- [`handlers.go`](./handlers.go) - HTTP 请求处理器
- [`models.go`](./models.go) - 数据结构
- [`store.go`](./store.go) - 数据持久化与业务逻辑

---

## 许可证

本项目按原样提供，仅供教育和内部使用目的。

---

## 其他文档

- [README.md](./README.md) - 英文版本文档
