package zap

import (
	"fmt"
	"go.uber.org/zap"
)

func InitZapLogger() (*zap.Logger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create zap logger: %w", err)
	}

	return logger, nil
}
