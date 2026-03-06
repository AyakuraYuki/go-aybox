package main

import (
	"math/rand"
	"time"

	"github.com/AyakuraYuki/go-aybox/log"
	"github.com/AyakuraYuki/go-aybox/log/console"
	"github.com/AyakuraYuki/go-aybox/statusbar"
)

func main() {
	bar := statusbar.New()
	bar.Start()
	defer bar.Stop()

	// Route all log output through bar.Writer() so that log lines and the
	// status bar rendering share the same mutex and never interleave.
	logger := log.New(log.WithWriters(console.New(console.WithWriter(bar.Writer()))))

	bar.SetTask("LoadConfig")
	simulateLogs(logger, "config", 8)

	bar.SetTask("ConnectDB")
	simulateLogs(logger, "database", 12)

	bar.SetTask("Migration (1/3)")
	simulateLogs(logger, "migration-users", 6)

	bar.SetTask("Migration (2/3)")
	simulateLogs(logger, "migration-orders", 10)

	bar.SetTask("Migration (3/3)")
	simulateLogs(logger, "migration-logs", 8)

	bar.SetTask("GenerateReport")
	simulateLogs(logger, "report", 15)

	bar.ClearTask()
	logger.Info().Msg("All tasks completed!")
}

func simulateLogs(logger *log.Logger, prefix string, count int) {
	levels := []func(name ...string) *log.Log{
		logger.Info,
		logger.Debug,
		logger.Warn,
	}
	for i := 1; i <= count; i++ {
		levels[rand.Intn(len(levels))](prefix).Bool("ok", true).Msgf("processing step %d/%d", i, count)
		time.Sleep(time.Duration(200+rand.Intn(400)) * time.Millisecond)
	}
}
