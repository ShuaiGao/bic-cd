package log

import (
	"bic-cd/internal/util"
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

// MidLogger 自定义Logger中间件
func MidLogger() gin.HandlerFunc {
	if gin.IsDebugging() {
		return debug()
	}
	return normal()
}

func debug() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		//解析Body
		body := parseBody(c)
		traceID := getTraceID(c) //生成请求trace_id
		c.Set(TraceKey, traceID)
		c.Set(TraceTimestamp, now.UnixMilli())
		c.Header(TraceKey, traceID)
		if "/health" != c.Request.URL.Path {
			infof(
				[]zapcore.Field{
					{Key: "client_ip", Type: zapcore.StringType, String: c.ClientIP()},
					{Key: "method", Type: zapcore.StringType, String: c.Request.Method},
					{Key: "url", Type: zapcore.StringerType, String: c.Request.URL.Path},
					{Key: TraceKey, Type: zapcore.StringType, String: traceID},
				},
				">> recv: body[ %s ] query[ %s ] size[%d]",
				string(body),
				c.Request.URL.RawQuery,
				len(body))
		}
		c.Next()
		userId := c.GetUint(util.UserID)
		if "/health" != c.Request.URL.Path {
			infof(
				[]zapcore.Field{
					{Key: "client_ip", Type: zapcore.StringType, String: c.ClientIP()},
					{Key: "method", Type: zapcore.StringType, String: c.Request.Method},
					{Key: "url", Type: zapcore.StringerType, String: c.Request.URL.Path},
					{Key: "user_id", Type: zapcore.Uint32Type, Integer: int64(userId)},
					{Key: TraceKey, Type: zapcore.StringType, String: traceID},
					{Key: "duration", Type: zapcore.Uint32Type, Integer: time.Since(now).Milliseconds()},
					{Key: "status", Type: zapcore.Uint32Type, Integer: int64(c.Writer.Status())},
				},
				"<< send size[%d]",
				c.Writer.Size())
		}
	}
}

func normal() gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		traceID := getTraceID(c) //生成请求trace_id
		c.Set(TraceKey, traceID)
		c.Set(TraceTimestamp, now.UnixMilli())
		c.Header(TraceKey, traceID)
		if "/health" != c.Request.URL.Path {
			infof(
				[]zapcore.Field{
					{Key: "client_ip", Type: zapcore.StringType, String: c.ClientIP()},
					{Key: "method", Type: zapcore.StringType, String: c.Request.Method},
					{Key: "url", Type: zapcore.StringerType, String: c.Request.URL.Path},
					{Key: TraceKey, Type: zapcore.StringType, String: traceID},
				},
				">> recv: query[ %s ]", c.Request.URL.RawQuery)
		}
		c.Next()
		userId := c.GetUint(util.UserID)
		if "/health" != c.Request.URL.Path {
			infof(
				[]zapcore.Field{
					{Key: "client_ip", Type: zapcore.StringType, String: c.ClientIP()},
					{Key: "method", Type: zapcore.StringType, String: c.Request.Method},
					{Key: "url", Type: zapcore.StringerType, String: c.Request.URL.Path},
					{Key: "user_id", Type: zapcore.Uint32Type, Integer: int64(userId)},
					{Key: TraceKey, Type: zapcore.StringType, String: traceID},
					{Key: "duration", Type: zapcore.Uint32Type, Integer: time.Since(now).Milliseconds()},
					{Key: "status", Type: zapcore.Uint32Type, Integer: int64(c.Writer.Status())},
				},
				"<< send size[%d]",
				c.Writer.Size())
		}
	}
}

func parseBody(c *gin.Context) []byte {
	var buf bytes.Buffer
	tee := io.TeeReader(c.Request.Body, &buf)
	body, err := io.ReadAll(tee)
	if err != nil {
		return nil
	}
	c.Request.Body = io.NopCloser(&buf)
	return body
}

func getTraceID(c *gin.Context) string {
	traceHeader := c.GetHeader(TraceKey)
	if traceHeader != "" {
		return traceHeader
	}
	return xid.New().String() // create a new x-trace-id
}
