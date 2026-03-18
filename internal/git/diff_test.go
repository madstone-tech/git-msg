package git_test

import (
	"context"
	"errors"
	"testing"

	"github.com/madstone0-0/git-msg/internal/git"
)

func TestFakeClient_StagedDiff_Empty(t *testing.T) {
	client := &git.FakeClient{DiffErr: git.ErrNoStagedChanges}
	ctx := context.Background()
	_, err := client.StagedDiff(ctx)
	if !errors.Is(err, git.ErrNoStagedChanges) {
		t.Errorf("expected ErrNoStagedChanges, got %v", err)
	}
}

func TestFakeClient_StagedDiff_NonEmpty(t *testing.T) {
	client := &git.FakeClient{DiffOut: "diff --git a/foo.go"}
	ctx := context.Background()
	out, err := client.StagedDiff(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if out == "" {
		t.Error("expected non-empty diff")
	}
}
