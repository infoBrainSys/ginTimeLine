package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
	"timeLineGin/pkg/config"
)

var instance *zap.SugaredLogger

func GetInstance() *zap.SugaredLogger {
	return instance
}

func Initialize(conf *config.Logger) {
	// 初始化日志
	var level zapcore.Level
	switch conf.Level {
	case -1:
		level = zap.DebugLevel
		break
	case 0:
		level = zap.InfoLevel
		break
	case 1:
		level = zap.WarnLevel
		break
	case 2:
		level = zap.ErrorLevel
		break
	case 3:
		level = zap.DPanicLevel
		break
	case 4:
		level = zap.PanicLevel
		break
	case 5:
		level = zap.FatalLevel
		break
	default:
		level = zap.InfoLevel
	}

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime = zapcore.ISO8601TimeEncoder   //转换编码的时间戳
	encoder.EncodeLevel = zapcore.CapitalLevelEncoder // 编码级别调整为大写的级别输出

	productConf := zapcore.NewCore(zapcore.NewJSONEncoder(encoder), loggerWrite(conf), level)
	logger := zap.New(productConf, zap.AddCaller()).Sugar() // AddCaller 将 Logger 配置为使用 zap 调用方的文件名、行号和函数名称对每条消息进行注释。

	instance = logger
}

func loggerWrite(conf *config.Logger) zapcore.WriteSyncer {
	fileLogger := &lumberjack.Logger{
		Filename:   conf.Path,
		MaxSize:    conf.MaxSize,
		MaxAge:     conf.MaxAge,
		MaxBackups: conf.MaxBackups,
		LocalTime:  true,
		Compress:   conf.Compress,
	}
	return zapcore.AddSync(fileLogger)
}

// GinLogger 将 gin 使用 zap log
func GinLogger(l *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		l.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}
