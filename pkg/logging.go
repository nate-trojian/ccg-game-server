package pkg

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _logger *zap.Logger = nil

func InitializeLogger(level string) error {
	var lvl zapcore.Level
	err := lvl.Set(level)
	if err != nil {
		return err
	}

	ec := zap.NewProductionEncoderConfig()
	ec.MessageKey = "m"
	ec.LevelKey = "l"
	ec.TimeKey = "t"
	ec.CallerKey = "c"
	ec.EncodeTime = zapcore.ISO8601TimeEncoder

	enc := zapcore.NewJSONEncoder(ec)
	ws := zapcore.Lock(os.Stdout)
	core := zapcore.NewCore(enc, ws, lvl)

	_logger = zap.New(core, zap.AddCaller())
	return nil
}

func NewLogger(name string) *zap.Logger {
	return _logger.With(zap.String("n", name))
}
