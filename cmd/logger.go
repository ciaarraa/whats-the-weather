package cmd

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)
var logger *zap.SugaredLogger

func toggleDebug(cmd *cobra.Command, args []string) {
	if debug {
		config := zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		mainLogger, err := config.Build()
		defer mainLogger.Sync()
		if err != nil {
			panic(err)
		}
		logger = mainLogger.Sugar()
	} else {
		config := zap.Config{
			Level:            zap.NewAtomicLevelAt(zapcore.ErrorLevel),
			Development:      false,
			Encoding:         "console",
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}

		mainLogger, err := config.Build()
		defer mainLogger.Sync()
		if err != nil {
			panic(err)
		}
		logger = mainLogger.Sugar()
	}
}
