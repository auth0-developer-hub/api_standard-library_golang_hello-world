package common

import (
	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger loads a global logger based on a koanf configuration
func InitLogger(settings map[string]string) {

	logConfig := zap.NewProductionConfig()
	logConfig.Sampling = nil

	// Log Level
	var logLevel zapcore.Level
	if err := logLevel.Set(settings["logger.level"]); err != nil {
		zap.S().Fatalw("Could not determine logger.level", "error", err)
	}
	logConfig.Level.SetLevel(logLevel)

	// Handle json logger encodings
	logConfig.Encoding = "json"
	logConfig.EncoderConfig = zapdriver.NewDevelopmentEncoderConfig()

	// Settings
	logConfig.Development = settings["logger.dev_mode"] == "true"
	logConfig.DisableCaller = settings["logger.disable_caller"] == "true"

	logConfig.OutputPaths = []string{settings["logger.outputpath"]}
	logConfig.ErrorOutputPaths = []string{settings["logger.error_outputpath"]}

	// Build the logger
	globalLogger, _ := logConfig.Build()
	zap.ReplaceGlobals(globalLogger)
}
