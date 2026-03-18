package hook_test

import (
	"testing"

	"github.com/madstone-tech/git-msg/internal/hook"
)

func TestShouldGenerate(t *testing.T) {
	cases := []struct {
		src  hook.SourceType
		want bool
	}{
		{hook.SourceNormal, true},
		{hook.SourceMessage, true},
		{hook.SourceTemplate, false},
		{hook.SourceMerge, false},
		{hook.SourceSquash, false},
		{hook.SourceCommit, false},
	}
	for _, c := range cases {
		got := hook.ShouldGenerate(c.src)
		if got != c.want {
			t.Errorf("ShouldGenerate(%q) = %v, want %v", c.src, got, c.want)
		}
	}
}
