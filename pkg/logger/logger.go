package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	//DefaultLevel the default level
	DefaultLevel = zapcore.InfoLevel

	//DefaultTimeLayout the default time layout
	DefaultTimeLayout = time.RFC3339
)

var (
	zapLog     *zap.Logger
	sugar      *zap.SugaredLogger
	filePath   = ""
	serverName = ""
)

func NewLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{level: DefaultLevel, fields: make(map[string]string)}
	for _, f := range opts {
		f(opt)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,       // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 设置日志级别
	var atom zap.AtomicLevel
	if opt.debug {
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	config := zap.Config{
		Level:         atom,                                      // 日志级别
		Encoding:      "console",                                 // 输出格式 console 或 json
		EncoderConfig: encoderConfig,                             // 编码器配置
		InitialFields: map[string]interface{}{"服务器": serverName}, // 初始化字段，如：添加一个服务器名称
		// Sampling: &zap.SamplingConfig{
		// 	Initial:    10,   // 每条记录中的一条会被记录
		// 	Thereafter: 1000, // 随后每 1000 条记录中的一条会被记录
		// },
	}

	if opt.debug {
		config.Development = true
	} else {
		config.Development = false
	}

	//logDir := xsf_server.RunDirPrefix + "/log/"
	//exsits, _ := xsf_util.PathExists(logDir)
	//if !exsits {
	//	os.Mkdir(logDir, os.ModePerm)
	//}

	//var logfile string
	//if false {
	//	logfile = fmt.Sprintf("%s%s.log", logDir, xsf_server.ServerTag)
	//} else {
	//	xsf_server.LogIndex = 1
	//	for {
	//		logfile = fmt.Sprintf("%s%s-%s-%d.log", logDir, xsf_server.ServerTag,
	//			xsf_util.EP2Name(int(xsf_server.SID.Type)),
	//			xsf_server.LogIndex)
	//
	//		_, err := os.Stat(logfile)
	//		if err == nil {
	//			xsf_server.LogIndex++
	//		} else {
	//			break
	//		}
	//	}
	//
	//}

	if opt.disableConsole {
		config.OutputPaths = []string{"stdout", filePath}
		config.ErrorOutputPaths = []string{"stderr", filePath}
	} else {
		config.OutputPaths = []string{filePath}
		config.ErrorOutputPaths = []string{filePath}
	}

	{
		var outputFile string
		outputFile = fmt.Sprintf("%s%s.race.log", filePath, serverName)

		err := os.Setenv("GORACE", fmt.Sprintf("log_path=%s", outputFile))
		if err != nil {
			panic(fmt.Sprintf("Failed to set environment variable:%v", err))
		}
	}

	// 构建日志
	var err error
	zapLog, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("log 初始化失败: %v", err))
	}

	sugar = zapLog.Sugar()

	return zapLog, nil
}

func NewJSONLogger(opts ...Option) (*zap.Logger, error) {
	opt := &option{level: DefaultLevel, fields: make(map[string]string)}
	for _, f := range opts {
		f(opt)
	}

	//timeLayout := DefaultTimeLayout
	//if opt.timeLayout != "" {
	//	timeLayout = opt.timeLayout
	//}

	var (
		/*自定义时间格式*/
		customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		}
		/*自定义日志级别显示*/
		customLevelEncoder = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(level.CapitalString())
		}
		/*自定义代码路径、行号输出*/
		customCallerEncoder = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + caller.TrimmedPath() + "]")
		}
	)

	// similar to zap.NewProductionEncoderConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger", // used by logger.Named(key); optional; useless
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace", // use by zap.AddStacktrace; optional; useless

		EncodeDuration: zapcore.MillisDurationEncoder,
		//EncodeCaller:   zapcore.ShortCallerEncoder, // 全路径编码器
		//EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendString(t.Format(timeLayout))
		//},
		//EncodeLevel:   zapcore.LowercaseLevelEncoder, // 小写编码器
		//LineEnding:    zapcore.DefaultLineEnding,
		EncodeTime:       customTimeEncoder,
		EncodeCaller:     customCallerEncoder, // 全路径编码器
		EncodeLevel:      customLevelEncoder,
		LineEnding:       "\n",
		ConsoleSeparator: " ",
	}

	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// lowPriority usd by info\debug\warn
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl < zapcore.ErrorLevel
	})

	// highPriority usd by error\panic\fatal
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= opt.level && lvl >= zapcore.ErrorLevel
	})

	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	stderr := zapcore.Lock(os.Stderr) // lock for concurrent safe

	core := zapcore.NewTee()

	if !opt.disableConsole {
		//Dev环境,日志级别使用带颜色的标识
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
				zapcore.NewMultiWriteSyncer(stdout),
				lowPriority,
			),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig),
				zapcore.NewMultiWriteSyncer(stdout),
				highPriority,
			),
		)
	}

	if opt.file != nil {
		core = zapcore.NewTee(core,
			zapcore.NewCore(jsonEncoder,
				zapcore.AddSync(opt.file),
				zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
					return lvl >= opt.level
				}),
			),
		)
	}

	zapLog = zap.New(core,
		zap.AddCaller(),
		zap.ErrorOutput(stderr),
	)

	for key, value := range opt.fields {
		zapLog = zapLog.WithOptions(zap.Fields(zapcore.Field{Key: key, Type: zapcore.StringType, String: value}))
	}

	sugar = zapLog.Sugar()

	return zapLog, nil
}

func GetLogger() *zap.Logger {
	return zapLog
}

func Release() {
	zapLog.Sync()
}

func With(fields ...zap.Field) *zap.Logger {
	return zapLog.With(fields...)
}

func Debug(message string, fields ...zapcore.Field) {
	zapLog.Debug(message, fields...)
}

func Info(message string, fields ...zapcore.Field) {
	zapLog.Info(message, fields...)
}

func Warn(message string, fields ...zapcore.Field) {
	zapLog.Warn(message, fields...)
}

func Error(message string, fields ...zapcore.Field) {
	zapLog.Error(message, fields...)
}

func DPanic(message string, fields ...zapcore.Field) {
	zapLog.DPanic(message, fields...)
}

func Panic(message string, fields ...zapcore.Field) {
	zapLog.Panic(message, fields...)
}

func Fatal(message string, fields ...zapcore.Field) {
	zapLog.Fatal(message, fields...)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	sugar.DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args...)
}
