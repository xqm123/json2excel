package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"xorm.io/core"
)

// 通过 Logger 实现 xorm 的 ILogger interface
var _ core.ILogger = (*Logger)(nil)

type Logger struct {
	*zap.SugaredLogger
	level   zapcore.Level
	showSQL bool
}

var logger *Logger

func InitLogger(logpath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:  logpath, // ⽇志⽂件路径
		MaxSize:   15,      // megabytes
		LocalTime: true,
		Compress:  true, // 是否压缩 disabled by default
	}
	w := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), w, zap.NewAtomicLevelAt(level))
	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = &Logger{log.Sugar(), level, false}

	logger.Info("DefaultLogger init success")
	return log
}

func GetLogger() *Logger {
	return logger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}

// 增加以下四个方法来实现 xorm 的 core.ILogger interface
func (l *Logger) Level() core.LogLevel {
	switch l.level {
	case zap.DebugLevel:
		return core.LOG_DEBUG
	case zap.InfoLevel:
		return core.LOG_INFO
	case zap.WarnLevel:
		return core.LOG_WARNING
	case zap.ErrorLevel:
		return core.LOG_ERR
	default:
		return core.LOG_UNKNOWN
	}
}

func (l *Logger) SetLevel(level core.LogLevel) {
	// Nothing To Do
}

func (l *Logger) IsShowSQL() bool {
	if l.Level() < core.LOG_INFO {
		l.showSQL = true
	}
	return l.showSQL
}

func (l *Logger) ShowSQL(show ...bool) {
	if len(show) > 0 {
		l.showSQL = show[0]
	} else {
		l.showSQL = true
	}
}

// 向两个相邻的 arg 之间注入一个空格
func injectSpace(args ...interface{}) []interface{} {
	l := len(args)
	if l <= 1 {
		return args
	}

	l += l - 1
	spaced := make([]interface{}, l, l)
	for i := range spaced {
		if i%2 == 0 {
			spaced[i] = args[i/2]
		} else {
			spaced[i] = " "
		}
	}
	return spaced
}
