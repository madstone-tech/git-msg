package cmd_test

import (
	"context"
	"testing"

	"github.com/madstone-tech/git-msg/cmd"
	"github.com/madstone-tech/git-msg/internal/hook"
)

func TestInstallHook(t *testing.T) {
	mgr := &hook.FakeManager{}
	err := cmd.InstallHook(context.Background(), cmd.HookOptions{Global: false, Manager: mgr})
	if err != nil {
		t.Fatal(err)
	}
	if len(mgr.InstallCalls) != 1 || mgr.InstallCalls[0] != false {
		t.Error("Install not called with global=false")
	}
}

func TestInstallHook_Global(t *testing.T) {
	mgr := &hook.FakeManager{}
	err := cmd.InstallHook(context.Background(), cmd.HookOptions{Global: true, Manager: mgr})
	if err != nil {
		t.Fatal(err)
	}
	if len(mgr.InstallCalls) != 1 || mgr.InstallCalls[0] != true {
		t.Error("Install not called with global=true")
	}
}

func TestUninstallHook(t *testing.T) {
	mgr := &hook.FakeManager{}
	err := cmd.UninstallHook(context.Background(), cmd.HookOptions{Global: false, Manager: mgr})
	if err != nil {
		t.Fatal(err)
	}
	if len(mgr.UninstallCalls) != 1 {
		t.Error("Uninstall not called")
	}
}
