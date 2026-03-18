package git

import "context"

// FakeClient is a test double for Client.
type FakeClient struct {
	DiffOut      string
	DiffErr      error
	BranchOut    string
	BranchErr    error
	LogOut       string
	LogErr       error
	CommitErr    error
	ConfigOut    string
	ConfigErr    error
	CommittedMsg string
	RepoRootOut  string
	RepoRootErr  error
}

func (f *FakeClient) StagedDiff(ctx context.Context) (string, error) {
	return f.DiffOut, f.DiffErr
}
func (f *FakeClient) CurrentBranch(ctx context.Context) (string, error) {
	return f.BranchOut, f.BranchErr
}
func (f *FakeClient) RecentLog(ctx context.Context, n int) (string, error) {
	return f.LogOut, f.LogErr
}
func (f *FakeClient) Commit(ctx context.Context, msg string) error {
	f.CommittedMsg = msg
	return f.CommitErr
}
func (f *FakeClient) RunConfig(ctx context.Context, key string) (string, error) {
	return f.ConfigOut, f.ConfigErr
}
func (f *FakeClient) RepoRoot(ctx context.Context) (string, error) {
	return f.RepoRootOut, f.RepoRootErr
}
