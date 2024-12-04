package logger

import "go.uber.org/zap"

func NewLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	return logger
}
