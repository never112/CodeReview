package webhook

import (
	"context"
	"github.com/google/go-github/v76/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"testing"
)

func TestTest(t *testing.T) {
	log.Printf("testing ")
}

func NewCommentClient1(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return github.NewClient(tc)
}

func TestTest123(t *testing.T) {
	token := os.Getenv("GITHUB_TOKEN")
	log.Printf("token: %s", token)
	body := "xxxxxxxxxxxxxabcdsdcsddcsd"

	ctx := context.Background()

	comment := &github.IssueComment{Body: &body}

	createdComment, _, err := NewCommentClient(token).Issues.CreateComment(ctx, "never112", "CodeReview", 5, comment)
	log.Printf("%v", createdComment)
	log.Printf("%v", err)

}
