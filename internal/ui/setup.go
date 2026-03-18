package ui

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
)

// DefaultModelFor returns the canonical default model for a given provider.
func DefaultModelFor(provider string) string {
	switch provider {
	case "anthropic":
		return "claude-haiku-4-5"
	case "openai":
		return "gpt-4o-mini"
	case "gemini":
		return "gemini-1.5-flash"
	case "ollama":
		return "llama3"
	default:
		return ""
	}
}

// ListOllamaModels queries the local Ollama daemon and returns available model
// names. Returns nil if Ollama is not running or produces no output.
func ListOllamaModels() []string {
	out, err := exec.Command("ollama", "list").Output()
	if err != nil {
		return nil
	}
	var models []string
	for i, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if i == 0 {
			// Skip header row ("NAME  ID  SIZE  MODIFIED")
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// First field is the model name (e.g. "llama3:latest")
		name := strings.Fields(line)[0]
		models = append(models, name)
	}
	return models
}

// ErrWizardAborted is returned when the user cancels the setup wizard.
var ErrWizardAborted = errors.New("setup cancelled")

// SetupResult carries the values collected by the wizard.
// The caller is responsible for persisting config and storing credentials.
type SetupResult struct {
	Provider string
	Model    string
	// APIKey is empty when Provider == "ollama".
	APIKey string
}

// RunSetupWizard runs the first-run TUI wizard and returns the collected values.
// It performs no I/O against config files or the keychain — that is the
// caller's responsibility.
// Returns ErrWizardAborted if the user quits without completing setup.
func RunSetupWizard() (SetupResult, error) {
	fmt.Println()
	fmt.Println("  git-msg — first-run setup")
	fmt.Println("  ─────────────────────────")
	fmt.Println()

	provider, err := pickProvider()
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return SetupResult{}, ErrWizardAborted
		}
		return SetupResult{}, err
	}

	model, err := pickModel(provider)
	if err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return SetupResult{}, ErrWizardAborted
		}
		return SetupResult{}, err
	}

	var apiKey string
	if provider != "ollama" {
		apiKey, err = pickAPIKey(provider)
		if err != nil {
			if errors.Is(err, huh.ErrUserAborted) {
				return SetupResult{}, ErrWizardAborted
			}
			return SetupResult{}, err
		}
	}

	return SetupResult{Provider: provider, Model: model, APIKey: apiKey}, nil
}

func pickProvider() (string, error) {
	var provider string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select LLM provider").
				Options(
					huh.NewOption("Anthropic (Claude)", "anthropic"),
					huh.NewOption("OpenAI (GPT)", "openai"),
					huh.NewOption("Google Gemini", "gemini"),
					huh.NewOption("Ollama (local)", "ollama"),
				).
				Value(&provider),
		),
	).Run()
	return provider, err
}

func pickModel(provider string) (string, error) {
	if provider == "ollama" {
		return pickOllamaModel()
	}
	model := DefaultModelFor(provider)
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Model name (%s)", provider)).
				Description(fmt.Sprintf("Default: %s", model)).
				Value(&model),
		),
	).Run()
	return model, err
}

func pickAPIKey(provider string) (string, error) {
	var apiKey string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("%s API key", provider)).
				Description("Stored in system keychain — never written to disk").
				EchoMode(huh.EchoModePassword).
				Value(&apiKey),
		),
	).Run()
	return apiKey, err
}

func pickOllamaModel() (string, error) {
	models := ListOllamaModels()

	if len(models) == 0 {
		model := DefaultModelFor("ollama")
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Ollama model name").
					Description("Could not reach Ollama — enter a model name manually").
					Value(&model),
			),
		).Run()
		return model, err
	}

	opts := make([]huh.Option[string], len(models))
	for i, m := range models {
		opts[i] = huh.NewOption(m, m)
	}
	model := models[0]
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Ollama model").
				Description(fmt.Sprintf("%d model(s) found via `ollama list`", len(models))).
				Options(opts...).
				Value(&model),
		),
	).Run()
	return model, err
}
