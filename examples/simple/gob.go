//go:build ignore

package main

import (
	"log/slog"

	"github.com/phillezi/gob"
)

func init() {
	logger := slog.New(gob.NewPrettyHandler(nil))
	slog.SetDefault(logger)
}

func main() {
	gob.New(gob.WithDefaultTarget("all")).Add(
		"all",
		gob.Static(),
	).Add(
		"clean",
		gob.Clean(),
	).Run()
}
