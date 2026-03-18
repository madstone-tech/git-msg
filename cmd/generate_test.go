package cmd_test

import (
	"context"
	"os"
	"testing"

	"github.com/madstone-tech/git-msg/cmd"
	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/git"
	"github.com/madstone-tech/git-msg/internal/llm"
	"github.com/madstone-tech/git-msg/internal/prompt"
	"github.com/madstone-tech/git-msg/internal/ui"
)

func defaultOpts() cmd.GenerateOptions {
	cfg := config.DefaultConfig()
	return cmd.GenerateOptions{
		Git: &git.FakeClient{
			DiffOut:   "diff --git a/foo.go",
			BranchOut: "main",
			LogOut:    "abc123 prev commit",
		},
		LLM:     &llm.FakeProvider{Response: "feat: add feature"},
		Config:  &config.FakeStore{Cfg: cfg},
		Secrets: nil,
		Templates: func() *prompt.FakeTemplateStore {
			s := prompt.NewFakeTemplateStore()
			s.Templates["conventional"] = prompt.Template{
				Name:   "conventional",
				System: "sys {{ diff }}",
				User:   "user {{ branch }}",
			}
			return s
		}(),
		Cfg: &cfg,
	}
}

func TestRun_DryRun(t *testing.T) {
	opts := defaultOpts()
	opts.DryRun = true
	// capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	err := cmd.Run(context.Background(), opts)
	w.Close()
	os.Stdout = old
	buf := make([]byte, 1024)
	n, _ := r.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	out := string(buf[:n])
	if out == "" {
		t.Error("expected output from dry-run")
	}
}

func TestRun_NoStagedChanges(t *testing.T) {
	opts := defaultOpts()
	opts.Git = &git.FakeClient{DiffErr: git.ErrNoStagedChanges}
	err := cmd.Run(context.Background(), opts)
	if err == nil {
		t.Fatal("expected error for no staged changes")
	}
}

func TestRun_HookMode_SkipsNonGenerate(t *testing.T) {
	opts := defaultOpts()
	opts.HookMode = true
	opts.HookSource = "merge" // ShouldGenerate returns false
	err := cmd.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("expected nil for non-generate source, got %v", err)
	}
}

// Scenario 4: edit inline — user edits message and it's used for commit
func TestRun_ReviewEditInline(t *testing.T) {
	opts := defaultOpts()
	opts.ReviewFunc = func(msg string) (ui.ReviewResult, error) {
		return ui.ReviewResult{Action: ui.ActionEditInline, Message: "edited: " + msg}, nil
	}
	err := cmd.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// FakeClient.CommittedMsg should contain the edited message
	fc := opts.Git.(*git.FakeClient)
	if fc.CommittedMsg != "edited: feat: add feature" {
		t.Errorf("unexpected committed message: %q", fc.CommittedMsg)
	}
}

// Scenario 5: open editor — result of editor is committed
func TestRun_ReviewOpenEditor(t *testing.T) {
	opts := defaultOpts()
	opts.ReviewFunc = func(msg string) (ui.ReviewResult, error) {
		return ui.ReviewResult{Action: ui.ActionOpenEditor, Message: "editor-message"}, nil
	}
	err := cmd.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fc := opts.Git.(*git.FakeClient)
	if fc.CommittedMsg != "editor-message" {
		t.Errorf("unexpected committed message: %q", fc.CommittedMsg)
	}
}

// Scenario 6: abort — no commit, exit 0
func TestRun_ReviewAbort(t *testing.T) {
	opts := defaultOpts()
	opts.ReviewFunc = func(msg string) (ui.ReviewResult, error) {
		return ui.ReviewResult{Action: ui.ActionAbort}, nil
	}
	err := cmd.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("abort should return nil, got: %v", err)
	}
	fc := opts.Git.(*git.FakeClient)
	if fc.CommittedMsg != "" {
		t.Errorf("expected no commit on abort, got: %q", fc.CommittedMsg)
	}
}

// Scenario 7: hook-mode — confirmed message written to hook file, not committed
func TestRun_HookMode_WritesFile(t *testing.T) {
	f, _ := os.CreateTemp("", "hook-msg-*")
	f.Close()
	defer os.Remove(f.Name())

	opts := defaultOpts()
	opts.HookMode = true
	opts.HookSource = "" // SourceNormal — generation proceeds
	opts.HookMsgFile = f.Name()
	opts.ReviewFunc = func(msg string) (ui.ReviewResult, error) {
		return ui.ReviewResult{Action: ui.ActionConfirm, Message: msg}, nil
	}
	err := cmd.Run(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "feat: add feature" {
		t.Errorf("unexpected hook file content: %q", string(data))
	}
	// Also verify git.Commit was NOT called
	fc := opts.Git.(*git.FakeClient)
	if fc.CommittedMsg != "" {
		t.Errorf("expected no git commit in hook mode, got: %q", fc.CommittedMsg)
	}
}
