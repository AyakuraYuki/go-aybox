package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AyakuraYuki/go-aybox/log"
	"github.com/AyakuraYuki/go-aybox/statusbar"
)

var logger = log.New()

func main() {
	bar := statusbar.New()

	bar.Start()
	defer bar.Stop()

	bar.SetTask("LoadConfig")
	simulateLogs("config", 8)

	bar.SetTask("ConnectDB")
	simulateLogs("database", 12)

	bar.SetTask("Migration (1/3)")
	simulateLogs("migration-users", 6)

	bar.SetTask("Migration (2/3)")
	simulateLogs("migration-orders", 10)

	bar.SetTask("Migration (3/3)")
	simulateLogs("migration-logs", 8)

	bar.SetTask("GenerateReport")
	simulateLogs("report", 15)

	bar.ClearTask()
	fmt.Println("\n✅ All tasks completed!")
}

func simulateLogs(prefix string, count int) {
	levels := []string{"INFO", "DEBUG", "WARN"}
	for i := 1; i <= count; i++ {
		level := levels[rand.Intn(len(levels))]
		fmt.Printf("[%s] %s: processing step %d/%d ...\n",
			level, prefix, i, count)
		logger.Info().Msg("xxx")
		time.Sleep(time.Duration(200+rand.Intn(400)) * time.Millisecond)
	}
}
