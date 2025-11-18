# 使用Ubuntu作为基础镜像
FROM ubuntu:22.04

# 设置环境变量以避免交互式安装
ENV DEBIAN_FRONTEND=noninteractive
ENV TZ=UTC

# 更新包管理器并安装基础依赖
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    gnupg \
    lsb-release \
    software-properties-common \
    ca-certificates \
    build-essential \
    git \
    bash \
    && rm -rf /var/lib/apt/lists/*

# 安装Node.js (使用NodeSource仓库)
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash - \
    && apt-get install -y nodejs \
    && rm -rf /var/lib/apt/lists/*

# 安装GitHub CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg \
    && chmod go+r /usr/share/keyrings/githubcli-archive-keyring.gpg \
    && echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
    && apt-get update \
    && apt-get install -y gh \
    && rm -rf /var/lib/apt/lists/*

# 创建用户目录结构
RUN mkdir -p /root/.config/gh \
    && mkdir -p /root/.claude \
    && mkdir -p /root/.local/bin

# 安装Claude Code
RUN curl -fsSL https://claude.ai/install.sh | bash

# 设置Claude Code token (请替换为实际的token)
ENV CLAUDE_API_TOKEN="your_claude_api_token_here"

# 创建GitHub配置
RUN cat > /root/.config/gh/config.yml << 'EOF'
# GitHub CLI Configuration
git_protocol: https
prompt: enabled
editor:
pager: less

# 替换为你的GitHub用户名
user: your_github_username
EOF

# 设置Git配置
RUN git config --global user.name "Your Name" \
    && git config --global user.email "your.email@example.com"

# 设置GitHub认证 (需要运行时登录)
# RUN echo "your_github_token" | gh auth login --with-token

# 创建工作目录
WORKDIR /workspace

# 设置环境变量
ENV PATH="/root/.local/bin:$PATH"
ENV CLAUDE_CONFIG_DIR="/root/.claude"

# 创建启动脚本
RUN cat > /root/entrypoint.sh << 'EOF'
#!/bin/bash
echo "=== CodeReview Docker Environment ==="
echo "Tools available:"
echo "- Git: $(git --version)"
echo "- GitHub CLI: $(gh version)"
echo "- Node.js: $(node --version)"
echo "- Claude Code: $(claude --version)"
echo ""
echo "Claude Code API Token: ${CLAUDE_API_TOKEN:0:10}..."
echo ""
echo "To login to GitHub, run: gh auth login"
echo "To start working, run: cd /workspace"
echo ""
# 保持容器运行
exec "$@"
EOF

RUN chmod +x /root/entrypoint.sh

# 设置默认命令
ENTRYPOINT ["/root/entrypoint.sh"]
CMD ["/bin/bash"]