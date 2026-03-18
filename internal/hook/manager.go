package hook

// Manager manages the prepare-commit-msg git hook lifecycle.
type Manager interface {
	Install(global bool) error
	Uninstall(global bool) error
	IsInstalled(global bool) (bool, error)
}
