package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Server  ServerConfig
	GitHub  GitHubConfig
	Review  ReviewConfig
	WorkDir string
}

type ServerConfig struct {
	Port string
}

type GitHubConfig struct {
	WebhookSecret string
	Token         string
}

type ReviewConfig struct {
	Command string
	Timeout int // seconds
}

func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		GitHub: GitHubConfig{
			WebhookSecret: getEnv("GITHUB_WEBHOOK_SECRET", ""),
			Token:         getEnv("GITHUB_TOKEN", ""),
		},
		Review: ReviewConfig{
			Command: getEnv("REVIEW_COMMAND", "claude /review  review结果用中文展示,review 时候参考 /promote/review.md 文件 --dangerously-skip-permissions"),
			Timeout: getEnvAsInt("REVIEW_TIMEOUT", 300), // 5 minutes default
		},
		WorkDir: getEnv("WORK_DIR", "/tmp/code-review"),
	}

	// Validate required configuration
	if cfg.GitHub.Token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN must be set")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
