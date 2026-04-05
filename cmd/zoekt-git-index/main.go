package main

import (
	"os"

	"github.com/sourcegraph/zoekt/cmd/zoekt-git-index/app"
)

func main() {
	os.Exit(app.Main())
}
