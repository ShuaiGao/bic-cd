package log

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TraceTimestamp string = "x-trace-timestamp"
	TraceKey       string = "x-trace-id"
)

type consoleConfig struct {
	Enable       bool          // 是否开启
	IsJsonFormat bool          // 是否输出json格式
	Level        zapcore.Level // 日志等级
}

type jsonConfig struct {
	consoleConfig
	lumberjack.Logger
}
