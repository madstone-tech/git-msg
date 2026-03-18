package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Client wraps git subprocess calls. All methods use exec.CommandContext.
type Client interface {
	StagedDiff(ctx context.Context) (string, error)
	CurrentBranch(ctx context.Context) (string, error)
	RecentLog(ctx context.Context, n int) (string, error)
	Commit(ctx context.Context, msg string) error
	RunConfig(ctx context.Context, key string) (string, error)
	RepoRoot(ctx context.Context) (string, error)
}

// ExecClient is the production implementation using os/exec.
type ExecClient struct{}

func (e ExecClient) StagedDiff(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "diff", "--cached").Output()
	if err != nil {
		return "", err
	}
	result := strings.TrimSpace(string(out))
	if result == "" {
		return "", ErrNoStagedChanges
	}
	return result, nil
}

func (e ExecClient) CurrentBranch(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (e ExecClient) RecentLog(ctx context.Context, n int) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "log", "--oneline", fmt.Sprintf("-%d", n)).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (e ExecClient) Commit(ctx context.Context, msg string) error {
	return exec.CommandContext(ctx, "git", "commit", "-m", msg).Run()
}

func (e ExecClient) RunConfig(ctx context.Context, key string) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "config", "--get", key).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func (e ExecClient) RepoRoot(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", ErrNotGitRepo
	}
	return strings.TrimSpace(string(out)), nil
}
