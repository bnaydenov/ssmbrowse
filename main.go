package main

import (
	"github.com/bnaydenov/ssmbrowse/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func main() {

	buildData := map[string]interface{}{
		"version": version,
		"commit":  commit,
		"date":    date,
		"builtBy": builtBy,
	}

	// fmt.Printf("%s", buildData["version"])
	cmd.Entrypoint(buildData)
}
