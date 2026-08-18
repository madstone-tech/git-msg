package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/madstone-tech/git-msg/internal/config"
	"github.com/madstone-tech/git-msg/internal/git"
	"github.com/madstone-tech/git-msg/internal/hook"
	"github.com/madstone-tech/git-msg/internal/llm"
	"github.com/madstone-tech/git-msg/internal/prompt"
	"github.com/madstone-tech/git-msg/internal/secret"
	"github.com/madstone-tech/git-msg/internal/ui"
)

// GenerateOptions holds all inputs for the generate command.
// Every field that touches I/O or infrastructure must be injected — Run()
// contains no construction logic and performs no I/O directly.
type GenerateOptions struct {
	DryRun      bool
	HookMode    bool
	HookMsgFile string
	HookSource  string
	Provider    string // optional per-run override
	Template    string // optional per-run override

	// Injectable dependencies — all required (set by cobra wiring or tests).
	Git       git.Client
	LLM       llm.Provider // nil = constructed via newLLMProvider
	Config    config.Store
	Secrets   secret.SecretStore
	Templates prompt.TemplateStore
	Cfg       *config.Config // pre-loaded config from PersistentPreRunE

	// Injectable behaviours — nil values fall back to real implementations.
	SpinnerFunc func(label string) ui.StopFunc                // nil = ui.Spinner
	ReviewFunc  func(message string) (ui.ReviewResult, error) // nil = ui.Review
}

// Run executes the generate command. No cobra imports, no construction logic,
// no direct I/O (spinner and review are injected).
func Run(ctx context.Context, opts GenerateOptions) error {
	// 1. Resolve config.
	var cfg config.Config
	if opts.Cfg != nil {
		cfg = *opts.Cfg
	} else {
		var err error
		cfg, err = opts.Config.Load()
		if err != nil {
			return fmt.Errorf("could not load config\n  → run git-msg to complete first-run setup: %w", err)
		}
	}

	// 2. Hook source check — exit early for unsupported sources.
	if opts.HookMode {
		if !hook.ShouldGenerate(hook.SourceType(opts.HookSource)) {
			return nil
		}
	}

	// 3. Apply per-run overrides.
	if opts.Provider != "" {
		cfg.Provider.Name = opts.Provider
	}
	if opts.Template != "" {
		cfg.Prompt.Default = opts.Template
	}

	// 4. Resolve LLM provider.
	provider := opts.LLM
	if provider == nil {
		var err error
		provider, err = NewLLMProvider(cfg, opts.Secrets)
		if err != nil {
			return fmt.Errorf("could not initialise LLM provider\n  → %w", err)
		}
	}

	// 5. Collect git context.
	diff, err := opts.Git.StagedDiff(ctx)
	if err != nil {
		return fmt.Errorf("no staged changes\n  → stage your changes with: git add <files>")
	}
	branch, _ := opts.Git.CurrentBranch(ctx)
	log, _ := opts.Git.RecentLog(ctx, 5)

	// 6. Render prompt template.
	tmpl, err := opts.Templates.Get(cfg.Prompt.Default)
	if err != nil {
		return fmt.Errorf("template %q not found\n  → run: git-msg prompt list", cfg.Prompt.Default)
	}
	renderer := prompt.NewRenderer()
	system, user, err := renderer.Render(tmpl, prompt.TemplateVars{
		Diff: diff, Branch: branch, Log: log,
	})
	if err != nil {
		return fmt.Errorf("template rendering failed: %w", err)
	}

	// 7. Call LLM — spinner is injected so Run stays I/O-free.
	spinFn := opts.SpinnerFunc
	if spinFn == nil {
		spinFn = ui.Spinner
	}
	stop := spinFn("Generating...")
	rawMessage, err := provider.Generate(ctx, system, user)
	stop()
	if err != nil {
		return fmt.Errorf("LLM request failed\n  → %w", err)
	}
	rawMessage = llm.CleanResponse(rawMessage)
	if rawMessage == "" {
		return fmt.Errorf("LLM returned an empty message\n  → try again or switch provider")
	}

	// 8. Dry-run: print and return without review.
	if opts.DryRun {
		fmt.Println(rawMessage)
		return nil
	}

	// 9. Review — injected so Run stays testable without a real TUI.
	reviewFn := opts.ReviewFunc
	if reviewFn == nil {
		reviewFn = ui.Review
	}
	result, err := reviewFn(rawMessage)
	if err != nil {
		return err
	}
	if result.Action == ui.ActionAbort {
		return nil
	}

	// 10. Commit or write hook file.
	if opts.HookMode {
		return os.WriteFile(opts.HookMsgFile, []byte(result.Message), 0644)
	}
	return opts.Git.Commit(ctx, result.Message)
}
