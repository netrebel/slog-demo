package main

import (
	"context"
	"log"
	"os"
	"runtime/debug"

	"golang.org/x/exp/slog"
)

func main() {
	buildInfo, _ := debug.ReadBuildInfo()

	log.Println("Original log package")

	slog.Info("----- Printing with slog -------")
	slog.Debug("Debug message not shown. Default is INFO")
	slog.Info("Info message")
	slog.Warn("Warning message")
	slog.Error("Error message")

	log.Println("------ Printing with slog JSONHandler ---------")

	/*
		You can also create your own Logger instance through the slog.New() method.
		It accepts a non-nil Handler interface which determines how the logs are formatted and where they are written to.
		Here's an example that uses the built-in JSONHandler type to format log records as JSON and send them to the standard output:
	*/

	//To customize DEBUG
	/* opts := slog.HandlerOptions{
		Level: slog.LevelDebug,
	}
	logger := slog.New(opts.NewJSONHandler(os.Stdout))*/

	logger := slog.New(slog.NewJSONHandler(os.Stdout)).With(
		slog.Group("program_info", // using a child logger to add the same attribute to all logs
			slog.Int("pid", os.Getpid()),
			slog.String("go_version", buildInfo.GoVersion),
		))
	logger.Debug("Debug message not shown. Default is INFO")
	logger.Info("Info message")
	logger.Warn("Warning message")
	logger.Error("Error message")

	/*
		Note that the SetDefault() method also updates the default logger used by the log package so that existing applications using log.Printf()
		and related methods can switch to structured logging:
	*/
	slog.SetDefault(logger)
	log.Println("Hello from old logger is now JSON")

	// One of the key advantages of logging in a structured format is the ability to add arbitrary attributes to log records in the form of key/value pairs.
	logger.Info(
		"incoming request",
		"method", "GET",
		"time_taken_ms", 158,
		"path", "/hello/world?q=search",
		"status", 200,
		"user_agent", "Googlebot/2.1 (+http://www.google.com/bot.html)",
	)

	// Property misalignment could lead to bad entries being. To prevent such mistakes, it's best to use strongly typed contextual attributes:
	logger.Info(
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
		),
	)

	// To guarantee type safety when adding contextual attributes to your records, you must use the LogAttrs() method
	logger.LogAttrs(
		context.Background(),
		slog.LevelInfo,
		"incoming request",
		slog.String("method", "GET"),
		slog.Int("time_taken_ms", 158),
		slog.String("path", "/hello/world?q=search"),
		slog.Int("status", 200),
		slog.String(
			"user_agent",
			"Googlebot/2.1 (+http://www.google.com/bot.html)",
		),
	)

}
