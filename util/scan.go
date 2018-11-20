package util

import (
	"github.com/chennqqi/gshark/util/githubsearch"

	"github.com/urfave/cli"

	"github.com/chennqqi/gshark/logger"
	"time"
)

func Scan(ctx *cli.Context) {
	// seconds
	var Interval time.Duration = 900

	if ctx.IsSet("time") {
		Interval = time.Duration(ctx.Int("time"))
	}

	logger.Log.Println("scan github code ")
	// use go keyword or not
	githubsearch.ScheduleTasks(Interval)
}
