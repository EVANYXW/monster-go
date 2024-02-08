/**
 * @api post logs.
 *
 * User: yunshengzhu
 * Date: 2022/5/4
 * Time: 1:14 PM
 */
package logs

import (
	"fmt"
	"github.com/sirupsen/logrus"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"net/url"
	"runtime"
	"strings"
)

type Logger struct {
	*logrus.Logger
	ServerName string
}

var Log *Logger

type Level logrus.Level

type hook struct {
	FileName string //输出日志的代码文件名称
	Line     string //打印日志的行
	Skip     int
	levels   []logrus.Level
}

type logger struct {
	filePath     string
	level        uint32
	format       string
	prettyPrint  bool
	isColor      bool
	maxSize      int
	maxBackups   int
	maxAge       int
	compress     bool
	esAddress    []string
	esUser       string
	esPassword   string
	serverName   string
	zsUrl        string
	reportCaller bool
}

type loggerOptions func(logger2 *logger)

func WithEsPassword(esPassword string) loggerOptions {
	return func(logs *logger) {
		logs.esPassword = esPassword
	}
}

func WithZsUrl(url string) loggerOptions {
	return func(logs *logger) {
		logs.zsUrl = url
	}
}

func WithEsUser(esUser string) loggerOptions {
	return func(logs *logger) {
		logs.esUser = esUser
	}
}

func WithEsAddress(esAddress []string) loggerOptions {
	return func(logs *logger) {
		logs.esAddress = esAddress
	}
}

func WithFilePath(filePath string) loggerOptions {
	return func(logs *logger) {
		logs.filePath = filePath
	}
}

func WithLevel(level uint32) loggerOptions {
	return func(logs *logger) {
		logs.level = level
	}
}

func WithFormat(format string) loggerOptions {
	return func(logs *logger) {
		logs.format = format
	}
}

func WithPrettyPrint(prettyPrint bool) loggerOptions {
	return func(logs *logger) {
		logs.prettyPrint = prettyPrint
	}
}

func WithIsColor(isColor bool) loggerOptions {
	return func(logs *logger) {
		logs.isColor = isColor
	}
}

func WithMaxSize(maxSize int) loggerOptions {
	return func(logs *logger) {
		logs.maxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) loggerOptions {
	return func(logs *logger) {
		logs.maxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) loggerOptions {
	return func(logs *logger) {
		logs.maxAge = maxAge
	}
}

func WithCompress(compress bool) loggerOptions {
	return func(logs *logger) {
		logs.compress = compress
	}
}

func WithServerName(serverName string) loggerOptions {
	return func(logs *logger) {
		logs.serverName = serverName
	}
}

func WithReportCaller(flag bool) loggerOptions {
	return func(logs *logger) {
		logs.reportCaller = flag
	}
}

func NewLogger(options ...loggerOptions) *Logger {
	logs := &logger{
		filePath:    "./log/logs.log",
		level:       5,
		format:      "json",
		prettyPrint: true,
		isColor:     false,
		maxSize:     100,
		maxBackups:  1000,
		maxAge:      3650,
		compress:    false,
	}

	for _, option := range options {
		option(logs)
	}

	logger := &Logger{}
	switch logs.format {
	case "text":
		logger.Logger = newLogger(logrus.Level(logs.level), &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05", ForceColors: logs.isColor}, NewHook())
	case "json":
		logger.Logger = newLogger(logrus.Level(logs.level), &logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05", PrettyPrint: logs.prettyPrint}, NewHook())
	default:
		logger.Logger = newLogger(logrus.Level(logs.level), &logrus.TextFormatter{FullTimestamp: true, TimestampFormat: "2006-01-02 15:04:05", ForceColors: logs.isColor}, NewHook())
	}

	loggerJack := &lumberjack.Logger{
		Filename:   logs.filePath,
		MaxSize:    logs.maxSize,    // 日志文件大小，单位是 MB
		MaxBackups: logs.maxBackups, // 最大过期日志保留个数
		MaxAge:     logs.maxAge,     // 保留过期文件最大时间，单位 天
		Compress:   logs.compress,   // 是否压缩日志，默认是不压缩。这里设置为true，压缩日志
	}
	logger.SetOutput(loggerJack) // logrus 设置日志的输出方式

	logger.SetReportCaller(logs.reportCaller)

	if len(logs.esAddress) > 0 {
		esh := newEsHook(logs.esAddress, logs.esUser, logs.esPassword)
		logger.AddHook(esh)
	}

	if logs.zsUrl != "" {
		u, err := url.Parse(logs.zsUrl)
		if err != nil {
			fmt.Println("err:", err)
		}
		zsh := newCsHook(u.Scheme+"://"+u.Host, u.Query().Get("user"), u.Query().Get("password"), logs.serverName)
		logger.AddHook(zsh)
	}

	logger.ServerName = logs.serverName
	Log = logger
	return logger
}

// 实现 logrus.Hook 接口
func (hook *hook) Fire(entry *logrus.Entry) error {
	fileName, line := findCaller(hook.Skip)
	entry.Data[hook.FileName] = fileName
	entry.Data[hook.Line] = line
	return nil
}

// 实现 logrus.Hook 接口
func (hook *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 自定义hook
func NewHook(levels ...logrus.Level) logrus.Hook {
	hook := hook{
		FileName: "filePath",
		Line:     "line",
		Skip:     5,
		levels:   levels,
	}
	if len(hook.levels) == 0 {
		hook.levels = logrus.AllLevels
	}
	return &hook
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	//fmt.Println("getCaller", pc, file, line, ok)
	if !ok {
		return "", 0
	}
	n := 0
	//获取执行代码的文件名
	for i := len(file) - 1; i > 0; i-- {
		if string(file[i]) == "/" {
			n++
			if n >= 2 {
				//fmt.Println(n >= 2, file)
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}

func findCaller(skip int) (string, int) {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return file, line
}

// 自定义logger
func newLogger(level logrus.Level, format logrus.Formatter, hook logrus.Hook) *logrus.Logger {
	log := logrus.New()
	log.Level = level
	log.SetFormatter(format)
	log.Hooks.Add(hook)
	return log
}
