package main

import (
	"log/slog"
	"time"

	"github.com/phillezi/gob"
)

var version = "-"

func init() {
	logger := slog.New(gob.NewPrettyHandler(nil))
	slog.SetDefault(logger)
}

func main() {
	_ = gob.New(gob.WithDefaultTarget("all")).
		Add(
			"all",
			gob.Static(),
		).Add(
		"clean",
		gob.Clean(),
	).Add(
		"test",
		gob.Test(
			gob.WithRace(),
			gob.WithCoverage(),
			gob.WithTimeout(30*time.Second),
			gob.WithVerbose(),
		),
	).Run()
	//.Matrix(
	//	gob.PopularOSes,
	//	gob.PopularArches,
	//)).Run()

	//b := gob.New(gob.WithExitOnError(false))
	//b.Add("all", gob.Static())
	//if err := b.Run(); err != nil {
	//	slog.Default().Error("Failed to build", "error", err)
	//	os.Exit(1)
	//}
}
