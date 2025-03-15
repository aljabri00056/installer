package logger

import (
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	logger *log.Logger
	level  Level = INFO
)

func init() {
	logger = log.New(os.Stdout, "", log.LstdFlags)
}

func SetLevel(l Level) {
	level = l
}

func Debug(format string, v ...any) {
	if level <= DEBUG {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

func Info(format string, v ...any) {
	if level <= INFO {
		logger.Printf("[INFO] "+format, v...)
	}
}

func Warn(format string, v ...any) {
	if level <= WARN {
		logger.Printf("[WARN] "+format, v...)
	}
}

func Error(format string, v ...any) {
	if level <= ERROR {
		logger.Printf("[ERROR] "+format, v...)
	}
}
