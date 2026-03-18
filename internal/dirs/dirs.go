// Package dirs resolves XDG-compliant configuration paths for git-msg.
//
// All config lives under $XDG_CONFIG_HOME/mdstn/git-msg/ (defaulting to
// ~/.config/mdstn/git-msg/ when XDG_CONFIG_HOME is unset), regardless of
// the host OS. This overrides os.UserConfigDir(), which returns the macOS
// ~/Library/Application Support on Darwin.
package dirs

import (
	"os"
	"path/filepath"
)

const (
	org = "mdstn"
	app = "git-msg"
)

// ConfigRoot returns the base directory for all git-msg configuration:
//
//	$XDG_CONFIG_HOME/mdstn/git-msg   (when XDG_CONFIG_HOME is set)
//	~/.config/mdstn/git-msg          (fallback)
func ConfigRoot() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, org, app), nil
}

// ConfigFile returns the path to the main config file.
func ConfigFile() (string, error) {
	root, err := ConfigRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "config.toml"), nil
}

// PromptsDir returns the directory for user-defined prompt templates.
func PromptsDir() (string, error) {
	root, err := ConfigRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "prompts"), nil
}
