package logger

import (
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/log/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New creates a new Zap logger
func New(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// NewWithOTEL creates a new Zap logger with OpenTelemetry integration
// Logs will be sent via OTLP along with traces and metrics
func NewWithOTEL(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
	}

	// Build the base logger for local output (console/file)
	baseLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	// Create OTEL zap core that exports logs via OTLP
	// This uses the global LoggerProvider set by telemetry.Initialize()
	otelCore := otelzap.NewCore("github.com/nmn3m/pulsar/backend",
		otelzap.WithLoggerProvider(global.GetLoggerProvider()),
	)

	// Combine the base core (for console output) with OTEL core (for OTLP export)
	// This ensures logs appear both locally and are exported to VictoriaLogs
	combinedCore := zapcore.NewTee(baseLogger.Core(), otelCore)

	// Create logger with combined core
	logger := zap.New(combinedCore,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return logger, nil
}
