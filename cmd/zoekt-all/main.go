// Command zoekt-all is a combined binary containing all zoekt tools as
// subcommands. Build with:
//
//	go build -o zoekt-all ./cmd/zoekt-all
//
// Usage:
//
//	zoekt-all <subcommand> [flags...]
//
// Subcommands use the zoekt tool name without the "zoekt-" prefix:
//
//	zoekt-all webserver -index /data/index
//	zoekt-all git-index -incremental /path/to/repo
//	zoekt-all mirror-github -org myorg -dest /data/repos
//
// Symlink dispatch is also supported: if the binary is invoked via a
// symlink named "zoekt-<subcmd>", it runs that subcommand directly.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	archiveindexapp "github.com/sourcegraph/zoekt/cmd/zoekt-archive-index/app"
	dynamicindexserverapp "github.com/sourcegraph/zoekt/cmd/zoekt-dynamic-indexserver/app"
	gitcloneapp "github.com/sourcegraph/zoekt/cmd/zoekt-git-clone/app"
	gitindexapp "github.com/sourcegraph/zoekt/cmd/zoekt-git-index/app"
	indexapp "github.com/sourcegraph/zoekt/cmd/zoekt-index/app"
	indexserverapp "github.com/sourcegraph/zoekt/cmd/zoekt-indexserver/app"
	mergeindexapp "github.com/sourcegraph/zoekt/cmd/zoekt-merge-index/app"
	mirrorbitbucketapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-bitbucket-server/app"
	mirrorgerritapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-gerrit/app"
	mirrorgiteaapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-gitea/app"
	mirrorgithubapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-github/app"
	mirrorgitilesapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-gitiles/app"
	mirrorgitlabapp "github.com/sourcegraph/zoekt/cmd/zoekt-mirror-gitlab/app"
	repoindexapp "github.com/sourcegraph/zoekt/cmd/zoekt-repo-index/app"
	sgindexserverapp "github.com/sourcegraph/zoekt/cmd/zoekt-sourcegraph-indexserver/app"
	testapp "github.com/sourcegraph/zoekt/cmd/zoekt-test/app"
	webserverapp "github.com/sourcegraph/zoekt/cmd/zoekt-webserver/app"
	searchapp "github.com/sourcegraph/zoekt/cmd/zoekt/app"
	"github.com/sourcegraph/zoekt/internal/cmdexec"
)

var commands = map[string]func(){
	"search":                  searchapp.Main,
	"index":                   indexapp.Main,
	"git-index":               func() { os.Exit(gitindexapp.Main()) },
	"git-clone":               gitcloneapp.Main,
	"webserver":               webserverapp.Main,
	"indexserver":             indexserverapp.Main,
	"dynamic-indexserver":     dynamicindexserverapp.Main,
	"sourcegraph-indexserver": sgindexserverapp.Main,
	"merge-index":             mergeindexapp.Main,
	"archive-index":           archiveindexapp.Main,
	"repo-index":              repoindexapp.Main,
	"test":                    testapp.Main,
	"mirror-github":           mirrorgithubapp.Main,
	"mirror-gitlab":           mirrorgitlabapp.Main,
	"mirror-gerrit":           mirrorgerritapp.Main,
	"mirror-gitea":            mirrorgiteaapp.Main,
	"mirror-gitiles":          mirrorgitilesapp.Main,
	"mirror-bitbucket-server": mirrorbitbucketapp.Main,
}

func main() {
	if exe, err := os.Executable(); err == nil {
		cmdexec.SelfPath = exe
	}

	// Busybox-style: if invoked as "zoekt-<subcmd>", dispatch directly.
	base := filepath.Base(os.Args[0])
	if sub := strings.TrimPrefix(base, "zoekt-"); sub != base && sub != "all" {
		if fn, ok := commands[sub]; ok {
			resetFlags()
			fn()
			return
		}
	}

	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	sub := strings.TrimPrefix(os.Args[1], "zoekt-")
	fn, ok := commands[sub]
	if !ok {
		fmt.Fprintf(os.Stderr, "unknown subcommand %q\n\n", sub)
		usage()
		os.Exit(2)
	}

	os.Args = os.Args[1:]
	os.Args[0] = "zoekt-" + sub

	resetFlags()
	fn()
}

// resetFlags replaces the default FlagSet so subcommands can register
// flags without colliding with flags from imported packages' init().
func resetFlags() {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s <subcommand> [flags...]\n\nSubcommands:\n", os.Args[0])
	var names []string
	for name := range commands {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %s\n", name)
	}
}
