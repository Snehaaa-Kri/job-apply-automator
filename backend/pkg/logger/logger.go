package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps zap.Logger so the rest of the app doesn't
// import zap directly. If we ever swap zap for something else,
// we only change this file. (Dependency Inversion)
//
// Node analogy: like creating a winston instance and exporting it.
// Everyone imports YOUR logger, not winston directly.
type Logger struct {
	*zap.Logger
}

// New creates a logger based on the environment.
// In development → human-readable colored output
// In production → JSON output (for log aggregators like Loki/Datadog)
//
// Node analogy: winston.createLogger({ format: ... })
func New(env string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	if env == "production" {
		// Production: JSON logs, only Info level and above
		zapLogger, err = zap.NewProduction()
	} else {
		// Development: pretty colored logs, Debug level and above
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapLogger, err = cfg.Build()
	}

	if err != nil {
		return nil, err
	}

	return &Logger{zapLogger}, nil
}

// Sync flushes any buffered log entries.
// Always call this with defer in main.go.
// Node analogy: like closing a write stream before process exits.
func (l *Logger) Sync() {
	_ = l.Logger.Sync()
}

func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}