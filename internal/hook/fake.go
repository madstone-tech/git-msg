package hook

// FakeManager is a test double for Manager.
type FakeManager struct {
	InstallErr      error
	UninstallErr    error
	InstalledResult bool
	InstalledErr    error
	InstallCalls    []bool
	UninstallCalls  []bool
}

func (f *FakeManager) Install(global bool) error {
	f.InstallCalls = append(f.InstallCalls, global)
	return f.InstallErr
}
func (f *FakeManager) Uninstall(global bool) error {
	f.UninstallCalls = append(f.UninstallCalls, global)
	return f.UninstallErr
}
func (f *FakeManager) IsInstalled(global bool) (bool, error) {
	return f.InstalledResult, f.InstalledErr
}
