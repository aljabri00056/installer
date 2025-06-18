package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/aljabri00056/installer/handler"
	"github.com/aljabri00056/installer/logger"
)

func main() {
	c := handler.GetConfigFromEnv()

	switch strings.ToLower(c.LogLevel) {
	case "debug":
		logger.SetLevel(logger.DEBUG)
	case "info":
		logger.SetLevel(logger.INFO)
	case "warn":
		logger.SetLevel(logger.WARN)
	case "error":
		logger.SetLevel(logger.ERROR)
	default:
		logger.SetLevel(logger.INFO)
	}

	lh := &handler.Handler{Config: c}
	addr := fmt.Sprintf(":%d", c.Port)
	l, err := net.Listen("tcp4", addr)
	if err != nil {
		logger.Error("failed to listen: %v", err)
		return
	}
	logger.Info("listening on %s...", addr)
	if err := http.Serve(l, lh); err != nil {
		logger.Error("server error: %v", err)
		return
	}
	logger.Info("exiting")
}
