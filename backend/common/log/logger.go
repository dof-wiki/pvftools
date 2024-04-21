package log

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/ctx"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	LevelDebug Level = zapcore.DebugLevel
	LevelInfo        = zapcore.InfoLevel
	LevelWarn        = zapcore.WarnLevel
	LevelError       = zapcore.ErrorLevel
	LevelPanic       = zapcore.PanicLevel
	LevelFatal       = zapcore.FatalLevel
)

const defaultLevel = LevelDebug
const (
	EnvLogLv = "LOGLV"
)

// 仅在发生错误以上的信息时打印堆栈
func traceOnError(level zapcore.Level) bool {
	return level > zap.ErrorLevel
}

// ZAP的OPTIONS
var (
	skipCaller  = zap.AddCallerSkip(1)
	levelEnable = zap.AddStacktrace(zap.LevelEnablerFunc(traceOnError))
)

func newDevLogger(level Level) *zap.Logger {
	devConfig := zap.NewDevelopmentConfig()
	devConfig.Level = zap.NewAtomicLevelAt(level)
	l, err := devConfig.Build(skipCaller, levelEnable)
	if err != nil {
		panic(err)
	}
	return l
}

var sugaredLogger *zap.SugaredLogger
var metric *zap.Logger

func initSuggar(level Level) {
	sugaredLogger = newDevLogger(level).Sugar()
}

func initMetric() {
	devConfig := zap.NewProductionConfig()
	devConfig.DisableCaller = true
	devConfig.DisableStacktrace = true
	devConfig.OutputPaths = []string{"stdout"}
	devConfig.ErrorOutputPaths = []string{"stdout"}
	l, err := devConfig.Build(skipCaller, levelEnable)
	if err != nil {
		panic(err)
	}
	metric = l
}

func LogDebug(format string, args ...interface{}) {
	sugaredLogger.Debugf(format, args...)
}

func LogInfo(format string, args ...interface{}) {
	sugaredLogger.Infof(format, args...)
}

func LogWarn(format string, args ...interface{}) {
	sugaredLogger.Warnf(format, args...)
}

func LogError(format string, args ...interface{}) {
	if ctx.Ctx != nil {
		runtime.EventsEmit(*ctx.Ctx, consts.EventProgressEnd, fmt.Sprintf(format, args...))
	}
	sugaredLogger.Errorf(format, args...)
}

func LogPanic(format string, args ...interface{}) {
	sugaredLogger.Panicf(format, args...)
}

func Metric(msg string, fields ...zap.Field) {
	metric.Info(msg, fields...)
}

func SetLevel(level Level) {
	initSuggar(level)
}

func init() {
	logLv := defaultLevel
	if os.Getenv(EnvLogLv) == "prod" {
		logLv = LevelInfo
	}
	initSuggar(logLv)
	initMetric()
}
