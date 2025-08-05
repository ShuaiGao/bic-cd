package log

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	gormlogger "gorm.io/gorm/logger"
)

//-------------------- for gorm --------------------

type gormLogger struct {
	level         gormlogger.LogLevel
	slowThreshold time.Duration
	cutoffSize    uint64 // 默认0 表示不截断
}

const (
	traceStr     = "%s [%.3fms] [rows:%v] %s"
	traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
	traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
)

var (
	gormSourceDir     string
	errRecordNotFound = errors.New("record not found")
)

func gorminit() {
	_, file, _, _ := runtime.Caller(0)
	// compatible solution to get gorm source directory with various operating systems
	gormSourceDir = regexp.MustCompile(`gorm\.go`).ReplaceAllString(file, "")
}

// ExportGormLogger must after glog init
//
//	slowThreshold 慢日志阈值
//	cutoffSize 设置截断字符长度，默认0不截断
//		如："123中文" => 截断长度为5
//		如："123！@#" => 截断长度为6
func ExportGormLogger(slowThreshold time.Duration, cutoffSize uint64) gormlogger.Interface {
	gorminit()
	gl := &gormLogger{
		level:         gormlogger.Warn,
		slowThreshold: slowThreshold,
		cutoffSize:    cutoffSize,
	}
	return gl
}

// LogMode 设置日志输出等级
func (gl *gormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	gl.level = level
	return gl
}

func (gl *gormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if gl.level >= gormlogger.Info {
		Infof(ctx, format, args...)
	}
}

func (gl *gormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if gl.level >= gormlogger.Warn {
		Warnf(ctx, format, args...)
	}
}

func (gl *gormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if gl.level >= gormlogger.Error {
		Errorf(ctx, format, args...)
	}
}

func (gl *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if gl.level <= gormlogger.Silent {
		return
	}
	elapsed := time.Since(begin)
	switch {
	case err != nil && gl.level >= gormlogger.Error && (!errors.Is(err, errRecordNotFound)):
		sql, rows := fc()
		sql = gl.cutoff(sql)
		if rows == -1 {
			gl.Error(ctx, traceErrStr, fileWithLineNum(), err, elapsed.Milliseconds(), "-", sql)
		} else {
			gl.Error(ctx, traceErrStr, fileWithLineNum(), err, elapsed.Milliseconds(), rows, sql)
		}
	case elapsed > gl.slowThreshold && gl.slowThreshold != 0 && gl.level >= gormlogger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", gl.slowThreshold)
		sql = gl.cutoff(sql)
		if rows == -1 {
			gl.Warn(ctx, traceWarnStr, fileWithLineNum(), slowLog, elapsed.Milliseconds(), "-", sql)
		} else {
			gl.Warn(ctx, traceWarnStr, fileWithLineNum(), slowLog, elapsed.Milliseconds(), rows, sql)
		}
	case gl.level == gormlogger.Info:
		sql, rows := fc()
		sql = gl.cutoff(sql)
		if rows == -1 {
			gl.Info(ctx, traceStr, fileWithLineNum(), elapsed.Milliseconds(), "-", sql)
		} else {
			gl.Info(ctx, traceStr, fileWithLineNum(), elapsed.Milliseconds(), rows, sql)
		}
	}
}

func (gl *gormLogger) cutoff(sql string) string {
	cutoff := gl.cutoffSize
	if cutoff == 0 || sql == "" {
		return sql
	} else {
		r := []rune(sql)
		if len(r) <= int(cutoff) {
			return sql
		}
		return string([]rune(sql)[:cutoff])
	}
}

// fileWithLineNum return the file name and line number of the current file
func fileWithLineNum() string {
	// the second caller usually from glog & gorm internal, so set i start from 3
	for i := 3; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, gormSourceDir) || strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}
