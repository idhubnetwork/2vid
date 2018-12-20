package logger

import (
	"time"

	"2vid/config"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log to used by gin
var Log = logrus.New()

// log file path and name
var logName = config.V.LogName
var logLevel = config.V.LogLevel

var logLevels = map[string]logrus.Level{
	"DEBUG": logrus.DebugLevel,
	"INFO":  logrus.InfoLevel,
	"WARN":  logrus.WarnLevel,
	"ERROR": logrus.ErrorLevel,
	"FATAL": logrus.FatalLevel,
}

func init() {
	lfsHook := newLfsHook(logLevel, 100000)
	Log.AddHook(lfsHook)
}

// newLfsHook set Log level and return a rotatelogs Hook
func newLfsHook(logLevel string, maxRemainCnt uint) logrus.Hook {
	writer_info, err := rotatelogs.New(
		logName+"_INFO.%Y-%m-%d-%H.log",
		// WithLinkName为最新的日志建立软连接，以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),

		// WithRotationTime设置日志分割的时间，这里设置为一小时分割一次
		rotatelogs.WithRotationTime(24*time.Hour),

		// WithMaxAge和WithRotationCount二者只能设置一个，
		// WithMaxAge设置文件清理前的最长保存时间，
		// WithRotationCount设置文件清理前最多保存的个数
		// rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)

	writer_debug, err := rotatelogs.New(
		logName+"_DEBUG.%Y-%m-%d-%H.log",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	writer_warn, err := rotatelogs.New(
		logName+"_WARN.%Y-%m-%d-%H.log",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(24*time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	writer_error, err := rotatelogs.New(
		logName+"_ERROR.%Y-%m-%d-%H.log",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*365),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	writer_fatal, err := rotatelogs.New(
		logName+"_FATAL.%Y-%m-%d-%H.log",
		rotatelogs.WithLinkName(logName),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*365),
		// rotatelogs.WithRotationCount(maxRemainCnt),
	)

	if err != nil {
		logrus.Errorf("config local file system for logger error: %v", err)
	}

	level, ok := logLevels[logLevel]

	if ok {
		Log.SetLevel(level)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}

	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer_debug,
		logrus.InfoLevel:  writer_info,
		logrus.WarnLevel:  writer_warn,
		logrus.ErrorLevel: writer_error,
		logrus.FatalLevel: writer_fatal,
	}, &logrus.TextFormatter{})

	return lfsHook
}

// Logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = path + "?" + raw
		}

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		Log.Infof("[2vid] %v | %3d | %13v | %15s | %-7s s% |",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
