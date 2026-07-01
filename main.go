// Package main is the entry point for the mbz MusicBrainz CLI.
package main

import (
	"os"

	"github.com/liuguancheng/musicbrainz-cli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
