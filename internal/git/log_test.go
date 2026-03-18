package git_test

import (
	"context"
	"strings"
	"testing"

	"github.com/madstone0-0/git-msg/internal/git"
)

func TestFakeClient_RecentLog(t *testing.T) {
	client := &git.FakeClient{LogOut: "abc123 first\ndef456 second\nghi789 third"}
	log, err := client.RecentLog(context.Background(), 3)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(log, "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 log lines, got %d", len(lines))
	}
}
