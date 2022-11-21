package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogCallback ...
type LogCallback func(entry zapcore.Entry) error

// InitLogger ...
func InitLogger(atom zap.AtomicLevel) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.999Z07:00")

	logger := zap.New(zapcore.NewCore(
		// zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
	_ = zap.ReplaceGlobals(logger)
}

// RegisterCallback ...
func RegisterCallback(cb LogCallback) {
	core := zap.L().Core()
	core = zapcore.RegisterHooks(core, cb)
	logger := zap.New(core)
	_ = zap.ReplaceGlobals(logger)
}
