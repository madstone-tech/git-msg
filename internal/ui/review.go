package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
)

// ReviewAction represents the user's choice in the review UI.
type ReviewAction int

const (
	ActionConfirm ReviewAction = iota
	ActionEditInline
	ActionOpenEditor
	ActionAbort
)

// ReviewResult holds the outcome of the review UI.
type ReviewResult struct {
	Action  ReviewAction
	Message string
}

// Review presents the commit message review TUI and returns the user's decision.
func Review(message string) (ReviewResult, error) {
	var action string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Generated commit message:").
				Description(message).
				Options(
					huh.NewOption("Use as-is", "confirm"),
					huh.NewOption("Edit inline", "edit"),
					huh.NewOption("Open $EDITOR", "editor"),
					huh.NewOption("Abort", "abort"),
				).
				Value(&action),
		),
	).Run()
	if err != nil {
		return ReviewResult{}, err
	}

	switch action {
	case "confirm":
		return ReviewResult{Action: ActionConfirm, Message: message}, nil

	case "edit":
		edited := message
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewText().
					Title("Edit commit message").
					Value(&edited),
			),
		).Run()
		if err != nil {
			return ReviewResult{}, err
		}
		return ReviewResult{Action: ActionEditInline, Message: strings.TrimSpace(edited)}, nil

	case "editor":
		result, err := OpenInEditor(message)
		if err != nil {
			// Fallback: inline edit with warning
			fmt.Fprintln(os.Stderr, "Warning: $EDITOR not set or failed, falling back to inline edit")
			edited := message
			_ = huh.NewForm(
				huh.NewGroup(
					huh.NewText().Title("Edit commit message").Value(&edited),
				),
			).Run()
			return ReviewResult{Action: ActionEditInline, Message: strings.TrimSpace(edited)}, nil
		}
		return ReviewResult{Action: ActionOpenEditor, Message: result}, nil

	case "abort":
		return ReviewResult{Action: ActionAbort}, nil
	}
	return ReviewResult{Action: ActionAbort}, nil
}
