package main

import (
	"fmt"
	"time"

	"github.com/AyakuraYuki/go-aybox/statusbar"
)

func main() {
	// create status bar with default options
	bar := statusbar.New()

	// or use customized options:
	// bar := statusbar.NewWithOptions(
	//     statusbar.WithRefreshRate(100),                // 100ms refresh rate
	//     statusbar.WithStyle(statusbar.MinimalStyle()), // ASCII style
	// )

	// start status bar in an asynchronous coroutine manner
	bar.Start()

	// ensure stop status bar before program exit
	defer bar.Stop()

	// ——— do some tasks ———

	bar.SetTask("initialization")
	doWork("loading config files...", 2*time.Second)

	bar.SetTask("connect to database")
	doWork("connecting to database...", 3*time.Second)

	bar.SetTask("migrate data (1/3)")
	doWork("migrating t_user...", 2*time.Second)

	bar.SetTask("migrate data (2/3)")
	doWork("migrating t_order...", 2*time.Second)

	bar.SetTask("migrate data (3/3)")
	doWork("migrating t_log...", 2*time.Second)

	bar.SetTask("summary")
	doWork("generating summary...", 3*time.Second)

	bar.ClearTask()
	fmt.Println("\n✅ All tasks done！")
}

func doWork(msg string, duration time.Duration) {
	fmt.Println(msg)
	time.Sleep(duration)
}
