package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func infof(field []zapcore.Field, format string, args ...interface{}) {
	logger.logger.With(field...).Sugar().Infof(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logger.C(ctx).Infof(format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	logger.C(ctx).Info(args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logger.C(ctx).Warnf(format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logger.C(ctx).Warn(args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	logger.C(ctx).Error(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	logger.C(ctx).Error(args...)
}

func XInfof(format string, args ...interface{}) {
	logger.sugaredLogger.Infof(format, args...)
}

func XInfo(args ...interface{}) {
	logger.sugaredLogger.Info(args...)
}

func XWarnf(format string, args ...interface{}) {
	logger.sugaredLogger.Warnf(format, args...)
}

func XWarn(args ...interface{}) {
	logger.sugaredLogger.Warn(args...)
}

func XError(args ...interface{}) {
	logger.sugaredLogger.Error(args...)
}

func (l *zapLogger) C(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		return l.sugaredLogger
	}
	fields := make([]interface{}, 0, 2+2*(len(extraCtxField)))
	if traceTimestamp := ctx.Value(TraceTimestamp); traceTimestamp != nil {
		if v, ok := traceTimestamp.(int64); ok {
			fields = append(fields, "duration", time.Now().UnixMilli()-v)
		}
	}
	for _, v := range extraCtxField {
		if cv := ctx.Value(v); cv != nil {
			fields = append(fields, v, cv)
		}
	}
	if len(fields) == 0 {
		return l.sugaredLogger
	}
	return l.sugaredLogger.With(fields)
}
