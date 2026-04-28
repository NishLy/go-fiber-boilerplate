package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger

func Init() {
	logger, err := zap.NewProduction(
		zap.WrapCore(func(core zapcore.Core) zapcore.Core {
			return zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				zapcore.AddSync(zapcore.Lock(zapcore.AddSync(zapcore.NewMultiWriteSyncer()))),
				zap.NewAtomicLevelAt(zap.DebugLevel),
			)
		}),
	)
	if err != nil {
		panic(err)
	}

	Log = logger
	Sugar = logger.Sugar()
}
