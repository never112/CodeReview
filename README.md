# 代码审查机器人

📖 **Documentation** | [English](README.en.md) | [简体中文](README.zh.md)

一个基于 Go 语言的智能代码审查机器人，专为 GitHub 拉取请求设计。该机器人集成了 Claude Code AI 审查功能，能够自动分析代码质量并提供专业的中文审查报告。

## 项目概述

该项目是一个现代化的自动化代码审查解决方案，通过 GitHub webhook 机制实现实时代码质量监控。系统采用模块化设计，支持多种代码审查工具，并默认集成 Claude Code AI 智能审查功能。

### 核心功能

- 🤖 **智能监听**：自动响应 GitHub PR 的 `opened`、`synchronize`、`reopened` 事件
- 📥 **自动克隆**：智能获取 PR 源码，支持分支切换和版本管理
- 🔍 **AI 驱动审查**：默认集成 Claude Code，支持多维度代码质量分析
- 💬 **专业评论**：生成结构化的中文审查报告，包含代码建议和改进方案
- 🔧 **安全认证**：采用 GitHub Personal Access Token 和 Webhook 签名双重验证
- 📊 **详细分析**：提供文件类型统计、执行时间监控和审查覆盖率分析

### 技术架构

```
CodeReview/
├── main.go          # 主程序入口，服务器启动和路由配置
├── config/          # 配置管理模块，环境变量处理
├── webhook/         # Webhook 事件处理器，GitHub API 交互
├── review/          # 代码审查执行引擎
├── git/            # Git 操作模块，仓库克隆和清理
├── prompt/         # AI 审查提示词模板
└── utils/          # 工具函数库
```

## 快速开始

### 环境要求

- **Go 1.24+**（项目使用最新 Go 版本）
- **Git** 版本控制工具
- **Claude Code CLI**（用于 AI 驱动的代码审查）
- **GitHub Personal Access Token**（需要仓库访问权限）

### 安装部署

1. **克隆项目**
```bash
git clone https://github.com/your-username/CodeReview.git
cd CodeReview
```

2. **依赖管理**
```bash
go mod tidy    # 下载并整理依赖
go mod download # 验证依赖完整性
```

3. **编译构建**
```bash
# 开发版本
go run main.go

# 生产版本
go build -o code-review .
./code-review
```

## 配置指南

### 必需配置项

```bash
# GitHub 个人访问令牌（需要 repo 权限）
export GITHUB_TOKEN="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# Webhook 签名密钥（用于验证请求来源）
export GITHUB_WEBHOOK_SECRET="your-webhook-secret-key"
```

### 可选配置项

```bash
# 服务器端口（默认：8080）
export PORT="8080"

# 自定义审查命令（已内置 Claude Code 配置）
export REVIEW_COMMAND="claude /review --system-prompt review时候参考prompt/review.md文件,review时参考prompt/baseReview.md文件,IMPORTANT: 请使用中文进行review结果展示 --dangerously-skip-permissions"

# 审查超时时间，单位秒（默认：300秒）
export REVIEW_TIMEOUT="300"

# 临时工作目录（默认：/tmp/code-review）
export WORK_DIR="/var/tmp/code-review"
```

### GitHub Personal Access Token 创建

1. 访问 [GitHub Settings](https://github.com/settings/tokens)
2. 选择 **Generate new token (classic)**
3. 配置权限范围：
   - ✅ `repo` - 访问私有仓库
   - ✅ `public_repo` - 访问公开仓库
   - ✅ `issues:write` - 创建评论权限
4. 生成并安全保存令牌

## GitHub Webhook 配置

### 创建 Webhook

1. 进入目标仓库 → **Settings** → **Webhooks**
2. 点击 **Add webhook**
3. 配置参数：
   - **Payload URL**: `http://your-server:8080/webhook`
   - **Content type**: `application/json`
   - **Secret**: 与 `GITHUB_WEBHOOK_SECRET` 相同
   - **Events**: 选择 **"Let me select individual events"**
   - 勾选 **"Pull requests"**

### 事件响应

机器人监听以下 PR 事件：
- `pull_request.opened` - 新 PR 创建时
- `pull_request.synchronize` - PR 新提交推送时
- `pull_request.reopened` - PR 重新开启时

## AI 审查配置

### 审查提示词模板

项目提供专业的 AI 审查提示词：

**`prompt/review.md`** - 代码审查标准
- 环境变量获取规范
- 默认值设置要求
- 注释说明规范

**`prompt/baseReview.md`** - 多维度审查策略
- 代码质量审查
- 性能优化建议
- 文档准确性检查
- 安全性扫描

### 自定义审查命令

**使用 Claude Code（推荐）**：
```bash
export REVIEW_COMMAND="claude /review --system-prompt $(cat prompt/review.md) --dangerously-skip-permissions"
```

**集成多种工具**：
```bash
export REVIEW_COMMAND="eslint . && go vet ./... && mypy . && claude /review"
```

**使用自定义脚本**：
```bash
export REVIEW_COMMAND="./scripts/custom-review.sh"
```

## 审查报告示例

机器人会生成如下格式的专业审查报告：

```markdown
## 🤖 代码审查结果

### 📋 审查输出

#### 代码质量分析
- 发现 3 个潜在问题需要修复
- 建议 2 处性能优化点
- 代码结构良好，可读性强

#### 安全检查
- ✅ 未发现明显安全漏洞
- ✅ 输入验证完善
- ⚠️ 建议：增加错误处理机制

### ✅ 状态

审查成功完成，整体代码质量良好。

### 📁 文件分析

本次审查涵盖 12 个文件：

**文件类型分布：**
- `.go`: 6 个文件
- `.md`: 3 个文件
- `.yaml`: 2 个文件
- `.json`: 1 个文件

### 💡 改进建议

1. **main.go:45** - 建议增加错误日志记录
2. **config/config.go:32** - 推荐使用环境变量默认值
3. **webhook/handler.go:156** - 优化并发处理逻辑
```

## 系统监控

### 健康检查

```bash
curl http://localhost:8080/health
# 响应：{"status":"ok"}
```

### 日志系统

应用采用结构化 JSON 日志：

```json
{
  "level": "info",
  "msg": "Starting webhook server",
  "port": "8080",
  "time": "2024-01-15T10:30:00Z"
}
```

日志内容包括：
- Webhook 事件接收记录
- 仓库克隆操作状态
- 审查命令执行详情
- GitHub API 交互结果
- 系统错误和异常信息

## 安全最佳实践

### 生产环境部署

1. **HTTPS 加密**：使用反向代理（Nginx/Caddy）启用 SSL/TLS
2. **密钥管理**：使用密钥管理服务存储敏感信息
3. **权限控制**：限制 GitHub Token 权限范围
4. **速率限制**：防止 API 滥用和 DDoS 攻击
5. **定期更新**：保持依赖包最新版本

### Docker 部署示例

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o code-review .

FROM alpine:latest
RUN apk --no-cache add ca-certificates git
WORKDIR /root/
COPY --from=builder /app/code-review .
EXPOSE 8080
CMD ["./code-review"]
```

## 故障排除

### 常见问题诊断

**1. 签名验证失败**
```bash
# 检查环境变量
echo $GITHUB_WEBHOOK_SECRET

# 验证 webhook 配置
# 确保 GitHub 仓库中的 Secret 与环境变量一致
```

**2. 仓库克隆失败**
```bash
# 测试 GitHub Token 权限
curl -H "Authorization: token $GITHUB_TOKEN" \
     https://api.github.com/user/repos

# 检查仓库访问权限
# 确保 Token 有足够权限访问目标仓库
```

**3. 审查命令执行失败**
```bash
# 测试审查命令
cd /tmp/test-repo
claude /review --dangerously-skip-permissions

# 检查 Claude CLI 安装
claude --version
```

**4. 性能优化问题**
- 调整 `REVIEW_TIMEOUT` 参数
- 优化审查命令复杂度
- 考虑增加并发处理限制

### 调试模式

启用详细日志记录：

```bash
export LOG_LEVEL="debug"
go run main.go
```

## 贡献指南

欢迎提交 Issue 和 Pull Request！

### 开发流程

1. Fork 项目仓库
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 提交更改：`git commit -m "Add new feature"`
4. 推送分支：`git push origin feature/new-feature`
5. 创建 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 添加必要的单元测试
- 更新相关文档
- 确保所有测试通过

## 许可证

本项目采用 **MIT 许可证**，详见 [LICENSE](LICENSE) 文件。

---

## 技术支持

如遇问题或需要技术支持，请：
- 提交 [GitHub Issue](https://github.com/your-username/CodeReview/issues)
- 查看 [项目文档](https://github.com/your-username/CodeReview/wiki)
- 联系维护团队
</div>

---

**Language / 语言:**
<button onclick="showLanguage('en')">English</button> |
<button onclick="showLanguage('zh')">中文</button>

<script>
function showLanguage(lang) {
  const enContent = document.getElementById('en-content');
  const zhContent = document.getElementById('zh-content');

  if (lang === 'en') {
    enContent.style.display = 'block';
    zhContent.style.display = 'none';
  } else {
    enContent.style.display = 'none';
    zhContent.style.display = 'block';
  }
}

// Default to Chinese
showLanguage('zh');
</script>