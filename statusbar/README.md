# Go CLI StatusBar

A non-blocking terminal status bar that runs as an independent goroutine, displaying real-time information at the bottom of your terminal without interfering with your main program's output.

## Features

| Feature | Description |
|---------|-------------|
| Current Time | Live clock in `HH:MM:SS` format |
| Elapsed Time | Timer since `Start()`, displayed as `Xh XXm XXs` |
| Task Name | Register the currently executing function/method name at any time |
| Non-blocking | Runs in a separate goroutine — never blocks your main logic |
| Thread-safe | All public methods are safe to call from any goroutine |

## Preview

```
Loading config...
Connecting to database...
├──────────────────────────────────────────────────────────────────────┤
 🕐 14:23:05 │ ⏱ 2m31s │ ⚡ ConnectDatabase
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    "yourmodule/statusbar"
)

func main() {
    bar := statusbar.New()
    bar.Start()
    defer bar.Stop()

    bar.SetTask("Initialize")
    time.Sleep(2 * time.Second)

    bar.SetTask("ProcessData")
    time.Sleep(3 * time.Second)

    bar.ClearTask()
    fmt.Println("Done")
}
```

## API

```go
// Create an instance
bar := statusbar.New()
bar := statusbar.NewWithOptions(
    statusbar.WithRefreshRate(100),                 // refresh interval in ms (default: 200)
    statusbar.WithStyle(statusbar.MinimalStyle()),  // ASCII style (no emoji)
)

// Control
bar.Start()           // start in a goroutine (non-blocking)
bar.Stop()            // stop and restore terminal state

// Task registration
bar.SetTask("MyFunc") // display the current task name
bar.ClearTask()       // clear the task display
```

## Styles

Two built-in styles are available:

**DefaultStyle** — uses emoji icons (🕐, ⏱, ⚡) and box-drawing characters.

**MinimalStyle** — pure ASCII, works in any terminal:

```
|----------------------------------------------------------------------|
 [T] 14:23:05 | [E] 2m31s | [>] ConnectDatabase
```

## Dependencies

```
golang.org/x/term  — terminal width detection
```

## Installation

Copy the `statusbar/` package into your project, then run:

```bash
go mod tidy
```
