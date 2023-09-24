package main

import (
	"github.com/alecthomas/kong"
	"github.com/glesica/flowork/internal/app/cmd"
	"log/slog"
	"os"
)

var CLI struct {
	cmd.GlobalOptions
	Run *cmd.RunOptions `help:"Run a workflow" cmd:""`
}

func main() {
	ctx := kong.Parse(
		&CLI,
		kong.Name("flowork"),
		kong.Description("A simple runner for linear data workflows"),
		kong.ShortUsageOnError(),
	)

	configureLogging()

	switch ctx.Command() {
	case "run <workflow>":
		err := cmd.Run(CLI.Run, CLI.GlobalOptions)
		if err != nil {
			ctx.FatalIfErrorf(err)
		}
	}
}

func configureLogging() {
	var level slog.Level
	switch {
	case CLI.Debug:
		level = slog.LevelDebug
	case CLI.Verbose:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: CLI.Debug,
		Level:     level,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
