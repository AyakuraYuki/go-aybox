package main

import (
	"os"
	"time"

	"github.com/AyakuraYuki/go-aybox/log"
	"github.com/AyakuraYuki/go-aybox/log/console"
	"github.com/AyakuraYuki/go-aybox/statusbar"
)

func main() {
	bar := statusbar.New(
		statusbar.WithOutput(os.Stdout))
	log.Configure(log.WithWriters(console.New(console.WithWriter(bar.Writer()))))

	bar.Start()
	defer bar.Stop()

	log.Info().Str("foo", "bar").Msg("hello")

	time.Sleep(2 * time.Second)

	bar.SetMessage("go run")
	go func() {
		defer bar.SetMessage("go end")
		for i := 0; i < 100; i++ {
			log.Info().Msg("working")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	time.Sleep(2 * time.Second)

	log.Info().Msg("goodbye")

	time.Sleep(200 * time.Millisecond)
}
