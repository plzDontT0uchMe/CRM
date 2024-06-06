package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func init() {
	logger = slog.Default()
	/*switch config.GetConfig().Env {
	case "local":
		logger = slog.Default()
	case "prod":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}*/
}

func formatMessageWithArgs(level string, message string, args ...any) string {
	var pairs []string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			pairs = append(pairs, fmt.Sprintf("%v=%v", args[i], args[i+1]))
		}
	}
	return level + " " + message + " " + strings.Join(pairs, " ")
}

func CreateLog(level string, message string, args ...any) {
	loggerFile, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory for config: %v", err)
	}
	loggerFile += "\\..\\..\\internal\\logger\\logger.log"
	if _, err = os.Stat(loggerFile); err != nil {
		loggerFile, _ = os.Getwd()
		loggerFile += "/logger.log"
	}

	f, err := os.OpenFile(loggerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("", err)
	}
	defer f.Close()

	switch level {
	case "debug":
		logger.Debug(message, args...)
	case "info":
		logger.Info(message, args...)
	case "warn":
		logger.Warn(message, args...)
	case "error":
		logger.Error(message, args...)
	}

	logger := log.New(f, "", log.LstdFlags)
	logger.Println(formatMessageWithArgs(strings.ToUpper(level), message, args...))
}
