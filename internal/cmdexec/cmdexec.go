// Package cmdexec provides subprocess helpers for invoking zoekt tools.
// When running as part of the combined binary (SelfPath is set), zoekt-*
// subprocess calls are rewritten to self-invocations with subcommands.
package cmdexec

import (
	"context"
	"os/exec"
	"strings"
)

// SelfPath is set by the combined binary to its own executable path.
// When non-empty, ZoektCommand rewrites zoekt-* calls to self-invocations.
var SelfPath string

// ZoektCommand returns an exec.Cmd for invoking a zoekt tool.
func ZoektCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	if SelfPath != "" && strings.HasPrefix(name, "zoekt-") {
		sub := strings.TrimPrefix(name, "zoekt-")
		return exec.CommandContext(ctx, SelfPath, append([]string{sub}, args...)...)
	}
	return exec.CommandContext(ctx, name, args...)
}

// Command is like ZoektCommand without a context.
func Command(name string, args ...string) *exec.Cmd {
	return ZoektCommand(context.Background(), name, args...)
}
