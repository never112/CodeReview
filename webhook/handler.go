package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"code-review/config"
	"code-review/git"
	"code-review/review"

	"github.com/google/go-github/v76/github"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	config *config.Config
	client *github.Client
}

type PullRequestEvent struct {
	Action      string             `json:"action"`
	Number      int                `json:"number"`
	Repository  Repository         `json:"repository"`
	PullRequest PullRequestPayload `json:"pull_request"`
}

type Repository struct {
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	CloneURL string `json:"clone_url"`
	Owner    Owner  `json:"owner"`
}

type Owner struct {
	Login string `json:"login"`
}

type PullRequestPayload struct {
	Head PullRequestCommit `json:"head"`
	Base PullRequestCommit `json:"base"`
}

type PullRequestCommit struct {
	Ref string `json:"ref"`
	SHA string `json:"sha"`
}

func NewHandler(cfg *config.Config) *Handler {
	// Use personal access token authentication
	tc := &http.Client{}
	client := github.NewClient(tc)

	return &Handler{
		config: cfg,
		client: client,
	}
}

func (h *Handler) HandlePullRequest(w http.ResponseWriter, r *http.Request) {
	// Verify webhook signature
	if h.config.GitHub.WebhookSecret != "" {
		signature := r.Header.Get("X-Hub-Signature-256")
		if signature == "" {
			logrus.Error("Missing X-Hub-Signature-256 header")
			http.Error(w, "Missing signature", http.StatusUnauthorized)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logrus.WithError(err).Error("Failed to read request body")
			http.Error(w, "Failed to read body", http.StatusBadRequest)
			return
		}

		if !h.verifySignature(body, signature, h.config.GitHub.WebhookSecret) {
			logrus.Error("Invalid webhook signature")
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		// Reset body for further processing
		r.Body = io.NopCloser(bytes.NewBuffer(body))
	}

	// Parse the webhook payload
	var event PullRequestEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		logrus.WithError(err).Error("Failed to decode webhook payload")
		http.Error(w, "Failed to decode payload", http.StatusBadRequest)
		return
	}

	// Only handle opened and synchronized events (new commits)
	if event.Action != "opened" && event.Action != "synchronize" && event.Action != "" && event.Action != "reopened" {
		logrus.WithField("action", event.Action).Info("Ignoring PR event")
		w.WriteHeader(http.StatusOK)
		return
	}

	logrus.WithFields(logrus.Fields{
		"repository": event.Repository.FullName,
		"pr_number":  event.Number,
		"action":     event.Action,
	}).Info("Processing pull request event")

	// Process the PR in a goroutine
	go h.processPullRequest(event)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"received"}`))
}

func (h *Handler) verifySignature(body []byte, signatureHeader, secret string) bool {
	if len(signatureHeader) < 7 {
		return false
	}

	signature, err := hex.DecodeString(signatureHeader[7:]) // Remove "sha256="
	if err != nil {
		return false
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expectedSignature := mac.Sum(nil)

	return hmac.Equal(signature, expectedSignature)
}

func (h *Handler) processPullRequest(event PullRequestEvent) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(h.config.Review.Timeout)*time.Second)
	defer cancel()

	// Clone the repository
	repoDir, err := git.CloneRepository(ctx, event.Repository.CloneURL, event.PullRequest.Head.Ref, h.config.WorkDir)
	if err != nil {
		logrus.WithError(err).Error("Failed to clone repository")
		h.postErrorComment(ctx, event, fmt.Sprintf("Failed to clone repository: %v", err))
		return
	}

	// Perform code review
	reviewResult, err := review.PerformCodeReview(ctx, repoDir, h.config.Review.Command)
	if err != nil {
		logrus.WithError(err).Error("Failed to perform code review")
		h.postErrorComment(ctx, event, fmt.Sprintf("Failed to perform code review: %v", err))
		return
	}

	// Post review results to PR
	if err := h.postReviewComment(ctx, event, reviewResult); err != nil {
		logrus.WithError(err).Error("Failed to post review comment")
	}

	// Cleanup
	if err := git.CleanupRepository(repoDir); err != nil {
		logrus.WithError(err).Warn("Failed to cleanup repository")
	}
}
func NewCommentClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}
func (h *Handler) postReviewComment(ctx context.Context, event PullRequestEvent, reviewResult string) error {
	gev := os.Getenv("GITHUB_TEST_ENV")
	log.Printf("%v", gev)

	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("## 🤖 Code Review Results\n\n%s", reviewResult)),
	}

	_, _, err := h.client.Issues.CreateComment(ctx, event.Repository.Owner.Login, event.Repository.Name, event.Number, comment)

	createdComment, _, err := NewCommentClient(os.Getenv("GITHUB_TOKEN")).Issues.CreateComment(ctx, "never112", "CodeReview", 5, comment)
	log.Printf("%v", createdComment)
	if err != nil {

		return fmt.Errorf("failed to create comment: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"repository": event.Repository.FullName,
		"pr_number":  event.Number,
	}).Info("Successfully posted review comment")

	return nil
}

func (h *Handler) postErrorComment(ctx context.Context, event PullRequestEvent, errorMsg string) error {
	comment := &github.IssueComment{
		Body: github.String(fmt.Sprintf("## ⚠️ Code Review Failed\n\n%s", errorMsg)),
	}

	_, _, err := h.client.Issues.CreateComment(ctx, event.Repository.Owner.Login, event.Repository.Name, event.Number, comment)
	if err != nil {
		return fmt.Errorf("failed to create error comment: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"repository": event.Repository.FullName,
		"pr_number":  event.Number,
	}).Info("Successfully posted error comment")

	return nil
}
