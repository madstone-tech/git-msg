package ui

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// OpenInEditor writes content to a temp file, opens it in $EDITOR, and
// returns the saved result. Returns an error if $EDITOR is unset or the
// editor process fails.
func OpenInEditor(content string) (string, error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return "", fmt.Errorf("$EDITOR is not set")
	}
	f, err := os.CreateTemp("", "git-msg-*.txt")
	if err != nil {
		return "", err
	}
	name := f.Name()
	defer func() { _ = os.Remove(name) }()

	if _, err := f.WriteString(content); err != nil {
		_ = f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	cmd := exec.Command(editor, name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
