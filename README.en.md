# Code Review Bot

A Go-based webhook server that automatically performs code reviews on GitHub pull requests and posts results back as comments.

## Features

- 🤖 Listens to GitHub pull request webhook events
- 📥 Automatically clones repository code when PR is opened/synchronized
- 🔍 Executes configurable review commands (defaults to `/review`)
- 💬 Posts review results as formatted comments on pull requests
- 🔧 Uses GitHub Personal Access Token authentication
- 📊 Provides detailed review output with file analysis

## Setup

### Prerequisites

- Go 1.21 or later
- Git
- A code review tool or script that can be executed via `/review` command

### Installation

1. Clone this repository:
```bash
git clone <repository-url>
cd code-review
```

2. Install dependencies:
```bash
go mod tidy
go mod download
```

3. Build the application:
```bash
go build -o code-review .
```

### Configuration

The application can be configured using environment variables:

**Required Configuration**

- **`GITHUB_TOKEN`**: A GitHub personal access token with repo access
- **`GITHUB_WEBHOOK_SECRET`**: Your GitHub webhook secret (for webhook signature verification)

**Optional Configuration**

- **`PORT`**: Server port (default: 8080)
- **`REVIEW_COMMAND`**: Command to execute for code review (default: `/review`)
- **`REVIEW_TIMEOUT`**: Review timeout in seconds (default: 300)
- **`WORK_DIR`**: Working directory for cloning repositories (default: `/tmp/code-review`)

### Authentication

The application uses GitHub Personal Access Token authentication:

```bash
export GITHUB_TOKEN="your_github_personal_access_token_here"
export GITHUB_WEBHOOK_SECRET="your_webhook_secret_here"
```

To create a GitHub Personal Access Token:
1. Go to GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. Click "Generate new token (classic)"
3. Select scopes: `repo` (for private repositories) and `public_repo` (for public repositories)
4. Generate token and copy it securely

### Setting up GitHub Webhook

1. Go to your GitHub repository Settings → Webhooks
2. Click "Add webhook"
3. Set **Payload URL** to: `http://your-server-url:8080/webhook`
4. Set **Content type** to `application/json`
5. Set **Secret** to your webhook secret
6. Select **"Let me select individual events"**
7. Check **"Pull requests"**
8. Click "Add webhook"

### Running the Application

**Development**
```bash
go run main.go
```

**Production**
```bash
./code-review
```

The server will start on port 8080 (or your configured port).

### Custom Review Command

The bot executes the command specified in the `REVIEW_COMMAND` environment variable in the cloned repository directory. You can customize this to run any code review tool.

**Examples**

Using Claude Code (as requested in requirements):
```bash
export REVIEW_COMMAND="claude /review"
```

Using a custom script:
```bash
export REVIEW_COMMAND="./scripts/my-review.sh"
```

Using multiple tools:
```bash
export REVIEW_COMMAND="eslint . && go vet ./... && mypy ."
```

### Webhook Events

The bot responds to:
- `pull_request.opened` - When a new PR is created
- `pull_request.synchronize` - When new commits are pushed to an existing PR

### Example Review Output

The bot posts formatted comments like:

```markdown
## 🤖 Code Review Results

**Review Command:** `claude /review`
**Duration:** 2m15s
**Exit Code:** 0

### 📋 Review Output

```
Great code quality! No issues found.
```

### ✅ Status

Review completed successfully.

### 📁 Files Analyzed

15 files were included in this review.

**File Types:**
- `.go`: 8 files
- `.md`: 3 files
- `.yml`: 2 files
- `.json`: 2 files
```

### Logging

The application uses structured logging with logrus in JSON format. Logs include:
- Webhook events received
- Repository cloning operations
- Review execution details
- API interaction results
- Error information

### Security Considerations

- Always use HTTPS in production
- Keep your webhook secret secure
- Limit the scope of your GitHub token to only necessary permissions
- Consider implementing rate limiting for production use
- Regularly update dependencies

### Troubleshooting

**Common Issues**

1. **"Invalid signature" error**
   - Check that `GITHUB_WEBHOOK_SECRET` matches your webhook secret in GitHub
   - Ensure you're using the correct secret in your webhook configuration

2. **"Failed to clone repository" error**
   - Check that your GitHub token has proper repository access
   - Verify the repository is public or your token has access to private repos

3. **"Failed to perform code review" error**
   - Ensure your review command is available in the system PATH
   - Check that the review command can execute in the repository directory
   - Review the timeout setting if reviews take a long time

**Health Check**

The application provides a health check endpoint:
```
GET /health
```

Returns: `{"status":"ok"}`

## License

This project is licensed under the MIT License.