# 代码审查机器人

一个基于 Go 语言的 webhook 服务器，可以自动对 GitHub 拉取请求进行代码审查，并将审查结果作为评论发布到 PR 上。

## 功能特点

- 🤖 监听 GitHub 拉取请求的 webhook 事件
- 📥 当 PR 创建或更新时自动克隆仓库代码
- 🔍 执行可配置的审查命令（默认为 `/review`）
- 💬 将审查结果作为格式化评论发布到拉取请求上
- 🔧 使用 GitHub 个人访问令牌进行身份验证
- 📊 提供详细的审查输出和文件分析

## 安装设置

### 环境要求

- Go 1.21 或更高版本
- Git
- 可以通过 `/review` 命令执行的代码审查工具或脚本

### 安装步骤

1. 克隆仓库：
```bash
git clone <repository-url>
cd code-review
```

2. 安装依赖：
```bash
go mod tidy
go mod download
```

3. 构建应用：
```bash
go build -o code-review .
```

### 配置说明

应用可以通过环境变量进行配置：

**必需配置**

- **`GITHUB_TOKEN`**: 具有仓库访问权限的 GitHub 个人访问令牌
- **`GITHUB_WEBHOOK_SECRET`**: GitHub webhook 密钥（用于 webhook 签名验证）

**可选配置**

- **`PORT`**: 服务器端口（默认：8080）
- **`REVIEW_COMMAND`**: 执行代码审查的命令（默认：`/review`）
- **`REVIEW_TIMEOUT`**: 审查超时时间，单位秒（默认：300）
- **`WORK_DIR`**: 克隆仓库的工作目录（默认：`/tmp/code-review`）

### 身份验证

应用使用 GitHub 个人访问令牌进行身份验证：

```bash
export GITHUB_TOKEN="你的github个人访问令牌"
export GITHUB_WEBHOOK_SECRET="你的webhook密钥"
```

创建 GitHub 个人访问令牌的步骤：
1. 访问 GitHub 设置 → 开发者设置 → 个人访问令牌 → 令牌（经典）
2. 点击 "Generate new token (classic)"
3. 选择权限范围：`repo`（用于私有仓库）和 `public_repo`（用于公开仓库）
4. 生成令牌并妥善保存

### 设置 GitHub Webhook

1. 进入你的 GitHub 仓库设置 → Webhooks
2. 点击 "Add webhook"
3. 设置 **Payload URL** 为：`http://your-server-url:8080/webhook`
4. 设置 **Content type** 为 `application/json`
5. 设置 **Secret** 为你的 webhook 密钥
6. 选择 "Let me select individual events"
7. 勾选 "Pull requests"
8. 点击 "Add webhook"

### 运行应用

**开发环境**
```bash
go run main.go
```

**生产环境**
```bash
./code-review
```

服务器将在端口 8080（或你配置的端口）上启动。

### 自定义审查命令

机器人会在克隆的仓库目录中执行 `REVIEW_COMMAND` 环境变量指定的命令。你可以自定义此命令来运行任何代码审查工具。

**示例**

使用 Claude Code（根据需求要求）：
```bash
export REVIEW_COMMAND="claude /review"
```

使用自定义脚本：
```bash
export REVIEW_COMMAND="./scripts/my-review.sh"
```

使用多个工具：
```bash
export REVIEW_COMMAND="eslint . && go vet ./... && mypy ."
```

### Webhook 事件

机器人响应以下事件：
- `pull_request.opened` - 创建新的 PR 时
- `pull_request.synchronize` - 向现有 PR 推送新提交时

### 审查输出示例

机器人会发布如下格式的评论：

```markdown
## 🤖 代码审查结果

**审查命令：** `claude /review`
**执行时间：** 2分15秒
**退出代码：** 0

### 📋 审查输出

```
代码质量很好！未发现问题。
```

### ✅ 状态

审查成功完成。

### 📁 分析文件

本次审查包含 15 个文件。

**文件类型：**
- `.go`: 8 个文件
- `.md`: 3 个文件
- `.yml`: 2 个文件
- `.json`: 2 个文件
```

### 日志记录

应用使用 logrus 进行结构化日志记录，格式为 JSON。日志包括：
- 接收到的 webhook 事件
- 仓库克隆操作
- 审查执行详情
- API 交互结果
- 错误信息

### 安全注意事项

- 生产环境中始终使用 HTTPS
- 保护好你的 webhook 密钥
- 限制 GitHub 令牌的权限范围，只授予必要的权限
- 在生产环境中考虑实施速率限制
- 定期更新依赖项

### 故障排除

**常见问题**

1. **"Invalid signature" 错误**
   - 检查 `GITHUB_WEBHOOK_SECRET` 是否与 GitHub 中的 webhook 密钥匹配
   - 确保你在 webhook 配置中使用了正确的密钥

2. **"Failed to clone repository" 错误**
   - 检查你的 GitHub 令牌是否具有适当的仓库访问权限
   - 验证仓库是否为公开仓库或令牌具有访问私有仓库的权限

3. **"Failed to perform code review" 错误**
   - 确保你的审查命令在系统 PATH 中可用
   - 检查审查命令是否可以在仓库目录中执行
   - 如果审查时间较长，请检查超时设置

**健康检查**

应用提供健康检查端点：
```
GET /health
```

返回：`{"status":"ok"}`

## 许可证

本项目采用 MIT 许可证。