package main

import (
	"errors"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger) // flushes buffer, if any

	e := errors.New("this is an error")

	url := "testurl"
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

	sugar.Error(e)
	sugar.Errorf("message with error %v", e)
	sugar.Errorw("message with error", zap.Error(e))
}
