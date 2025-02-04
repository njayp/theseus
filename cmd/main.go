package main

import (
	"log/slog"

	"github.com/njayp/theseus/pkg/server"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	if err := server.NewServer().Start(8080); err != nil {
		panic(err)
	}
}
