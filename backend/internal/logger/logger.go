// internal/logger/logger.go - PERBAIKAN
package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger

func InitLogger(env string) {
	var config zap.Config
	
	if env == "production" {
		config = zap.NewProductionConfig()
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.OutputPaths = []string{"stdout", "/var/log/backend.log"}
		config.ErrorOutputPaths = []string{"stderr"}
		config.DisableStacktrace = true
	} else {
		config = zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	}

	// Add caller information
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var err error
	Log, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	
	Sugar = Log.Sugar()

	// Redirect standard log to zap
	zap.RedirectStdLog(Log)
	
	Log.Info("Logger initialized", 
		zap.String("environment", env),
		zap.String("level", config.Level.String()),
	)
}

// Helper functions untuk logging yang konsisten
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	if len(fields) > 0 {
		Log.Info(msg, fields...)
	} else {
		Log.Info(msg)
	}
}

func Error(ctx context.Context, msg string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	Log.Error(msg, allFields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}