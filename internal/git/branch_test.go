package git_test

import (
	"context"
	"testing"

	"github.com/madstone-tech/git-msg/internal/git"
)

func TestFakeClient_CurrentBranch(t *testing.T) {
	client := &git.FakeClient{BranchOut: "main"}
	branch, err := client.CurrentBranch(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if branch != "main" {
		t.Errorf("expected main, got %q", branch)
	}
}
