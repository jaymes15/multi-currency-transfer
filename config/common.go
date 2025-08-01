package config

import (
	"log/slog"
	"os"
)

var configurations Config
var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
}))
