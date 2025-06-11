package log

import (
	"bic-cd/pkg/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

type zapLogger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

var (
	logger        *zapLogger
	extraCtxField = []string{TraceKey, TraceTimestamp}
)

func Stop() {
	logger.logger.Sync()
}

func Setup() {
	Console := &consoleConfig{
		Enable:       false,
		IsJsonFormat: false,
		Level:        zapcore.InfoLevel,
	}
	Json := &jsonConfig{
		consoleConfig: consoleConfig{
			Enable:       true,
			IsJsonFormat: true,
			Level:        zapcore.InfoLevel,
		},
		Logger: lumberjack.Logger{
			Filename:   config.AppSetting.LogPath,
			MaxSize:    32,
			MaxBackups: 10,
			MaxAge:     30,
		},
	}
	if gin.IsDebugging() {
		Console.Enable = true
		Console.Level = zapcore.DebugLevel
	}
	logger = newZapLogger(Console, Json)
}

func newZapLogger(console *consoleConfig, json *jsonConfig) *zapLogger {
	cores := make([]zapcore.Core, 0)
	l := &zapLogger{}
	if console.Enable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(console.Level)
		writer := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(getEncoder(console.IsJsonFormat), writer, lvl)
		cores = append(cores, core)
	}
	if json.Enable {
		lvl := zap.NewAtomicLevel()
		lvl.SetLevel(json.Level)
		writer := zapcore.AddSync(&json.Logger)
		core := zapcore.NewCore(getEncoder(json.IsJsonFormat), writer, lvl)
		cores = append(cores, core)
	}
	combinedCore := zapcore.NewTee(cores...)
	opts := []zap.Option{zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(2)}
	l.logger = zap.New(combinedCore, opts...)
	l.sugaredLogger = l.logger.Sugar()
	return l
}

func getEncoder(isJSON bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.StacktraceKey = "stack"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if isJSON {
		return zapcore.NewJSONEncoder(encoderConfig)
	}
	encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}
