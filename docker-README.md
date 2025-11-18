# CodeReview Docker Environment

这个Docker环境包含了进行代码审查所需的所有工具。

## 包含的工具

- **Git** - 版本控制系统
- **GitHub CLI (gh)** - GitHub命令行工具
- **Bash** - Shell环境
- **Node.js** - JavaScript运行时
- **Claude Code** - Claude AI代码助手

## 快速开始

1. 复制环境变量文件：
```bash
cp .env.example .env
```

2. 编辑 `.env` 文件，填入你的token：
```
CLAUDE_API_TOKEN=your_actual_claude_token
GITHUB_TOKEN=your_actual_github_token
```

3. 构建并运行容器：
```bash
docker-compose up --build
```

## 使用方法

### 第一次设置

1. 进入容器：
```bash
docker-compose run --rm codereview bash
```

2. 登录GitHub（如果未设置token）：
```bash
gh auth login
```

3. 验证Claude Code：
```bash
claude --version
```

### 日常使用

启动容器并开始工作：
```bash
docker-compose up
```

或者一次性运行命令：
```bash
docker-compose run --rm codereview claude "review this code"
```

## 配置说明

### Claude Code Token
- 在 `.env` 文件中设置 `CLAUDE_API_TOKEN`
- 或者在Dockerfile中直接修改 `ENV CLAUDE_API_TOKEN`

### GitHub配置
- 在 `docker-compose.yml` 中设置GitHub用户信息
- 或者使用GitHub CLI进行登录：`gh auth login`

### Git配置
在Dockerfile中已预设基础配置，可按需修改：
```dockerfile
RUN git config --global user.name "Your Name" \
    && git config --global user.email "your.email@example.com"
```

## 数据持久化

- 工作目录挂载到 `/workspace`
- Claude Code缓存挂载到 `/root/.cache/claude`

## 故障排除

1. **Claude Code无法访问**：检查API token是否正确设置
2. **GitHub CLI认证失败**：运行 `gh auth logout` 然后重新登录
3. **权限问题**：确保工作目录权限正确