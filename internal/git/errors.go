package git

import "errors"

var (
	// ErrNoStagedChanges is returned when git diff --cached produces no output.
	ErrNoStagedChanges = errors.New("no staged changes")
	// ErrNotGitRepo is returned when the working directory is not a git repository.
	ErrNotGitRepo = errors.New("not a git repository")
)
