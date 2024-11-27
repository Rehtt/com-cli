package main

import (
	"log/slog"
	"os"

	"github.com/Rehtt/gocui"
)

var Version string

func main() {
	a, err := NewApp(&Info{
		Name: "com-cli",
		Ver:  Version,
	})
	if err != nil {
		slog.Error("NewApp", "error", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil && err != gocui.ErrQuit {
		slog.Error("Run", "error", err)
		os.Exit(1)
	}
	return
}
