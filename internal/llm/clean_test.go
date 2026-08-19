package llm

import (
	"testing"
)

func TestCleanResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no cleaning needed",
			input:    "feat(ui): add new button",
			expected: "feat(ui): add new button",
		},
		{
			name:     "markdown fences",
			input:    "```git\nfeat(ui): add new button\n```",
			expected: "feat(ui): add new button",
		},
		{
			name:     "markdown fences with text",
			input:    "Here is the message:\n```\nfeat(ui): add new button\n```",
			expected: "feat(ui): add new button",
		},
		{
			name:     "conversational filler",
			input:    "Here is your commit message: feat(ui): add new button",
			expected: "feat(ui): add new button",
		},
		{
			name:     "case insensitive filler",
			input:    "COMMIT MESSAGE: feat(ui): add new button",
			expected: "feat(ui): add new button",
		},
		{
			name:     "filler with colon and space",
			input:    "Suggested commit message:   feat(ui): add new button  ",
			expected: "feat(ui): add new button",
		},
		{
			name:     "complex leakage",
			input:    "Based on the diff, here is the suggested commit message:\n\n```\nfix(core): resolve race condition in buffer\n```\n\nI hope this helps!",
			expected: "fix(core): resolve race condition in buffer",
		},
		{
			name:     "empty input",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanResponse(tt.input)
			if got != tt.expected {
				t.Errorf("CleanResponse(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}
