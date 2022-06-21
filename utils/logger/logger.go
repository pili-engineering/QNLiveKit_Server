package logger

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

const RequestIDHeaderKey = "X-Reqid"
const LoggerCtxKey = "logger"

var (
	WithHook   bool
	LoggerHook Hook

	XReqIDConst = "X-ReqID"
	pid         = uint32(time.Now().UnixNano() % 4294967291)
)

var (
	StdLog *Logger
	MgoLog *Logger
)

func init() {
	// set stdlog and set level
	StdLog = NewEmptyLogger()
	MgoLog = NewEmptyLogger()
}

type Log interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Panic(args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Xput(logs []string)
	ReqID() string
}

type Logger struct {
	mu       sync.Mutex
	out      io.Writer
	LogEntry *log.Entry
	fields   map[string]interface{}
	reqId    string
}

// GetLoggerFromReq get logger from request
func GetLoggerFromReq(req *http.Request) *Logger {
	reqID := req.Header.Get(RequestIDHeaderKey)
	if reqID == "" {
		reqID = GenReqID()
	}
	return New(reqID)
}

// ReqLogger Get Log from ctx
func ReqLogger(ctx context.Context) *Logger {
	l := ctx.Value(LoggerCtxKey)
	if log, ok := l.(*Logger); ok {
		return log
	}

	// default use std logger
	return StdLog
}

func (logger *Logger) WithFields(fields map[string]interface{}) {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}

	logger.LogEntry = logger.LogEntry.WithFields(f)
}

func (logger *Logger) WithFieldsNewLogger(fields map[string]interface{}) (entry *log.Entry) {
	f := make(log.Fields)
	for k, v := range fields {
		f[k] = v
	}

	return logger.LogEntry.WithFields(f)
}

func (logger *Logger) SetLevel(level log.Level) {
	logger.mu.Lock()
	defer logger.mu.Unlock()
	logger.LogEntry.Logger.SetLevel(level)
}

// SetOutput sets the standard logger output.
func (logger *Logger) SetOutput(out io.Writer) {
	logger.mu.Lock()
	defer logger.mu.Unlock()

	logger.LogEntry.Logger.Out = out
}

func GenReqID() string {
	var b [12]byte
	binary.LittleEndian.PutUint32(b[:], pid)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

// New this function is recommended to create a new logger object
func New(o ...interface{}) *Logger {
	var reqID = ""
	if len(o) > 0 && o[0] != nil {
		a := o[0]

		switch a.(type) {
		case *Logger:
			return a.(*Logger)
		case *log.Entry:
			return NewLogger(a.(*log.Entry))
		case string:
			reqID = a.(string)
		}
	}

	if len(reqID) == 0 {
		reqID = GenReqID()
	}

	l := NewEmptyLoggerWithFields(map[string]interface{}{XReqIDConst: reqID})
	l.reqId = reqID
	return l
}

func NewLogger(entry *log.Entry) *Logger {
	l := &Logger{
		LogEntry: entry,
		fields:   make(map[string]interface{}),
	}
	SetFormat("json", l)
	return l
}

func NewEmptyLogger() *Logger {
	logger := log.New()
	l := &Logger{
		LogEntry: log.NewEntry(logger),
		fields:   make(map[string]interface{}),
	}
	SetFormat("json", l)
	return l
}

type LoggerFields map[string]interface{}

func NewEmptyLoggerWithFields(fields map[string]interface{}) *Logger {
	l := NewEmptyLogger()
	l.fields = fields
	SetFormat("json", l)
	return l
}

func DecorateLog(logger *log.Entry) *log.Entry {
	var (
		fileName, funcName string
	)
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		fileName = "???"
		funcName = "???"
		line = 0
	} else {
		funcName = runtime.FuncForPC(pc).Name()
		fileSlice := strings.Split(file, path.Dir(funcName))
		fileName = filepath.Join(path.Dir(funcName), fileSlice[len(fileSlice)-1]) + ":" + strconv.Itoa(line)
	}

	return logger.WithField("file", fileName).WithField("func", funcName)

}

// hook
type Hook interface {
	Levels() []log.Level
	Fire(*log.Entry) error
}

func (logger *Logger) AddHook(hook Hook) *Logger {
	if WithHook && hook != nil {
		logger.LogEntry.Logger.AddHook(hook)
		return logger
	}
	return logger
}

func (logger *Logger) AddDefaultHook() *Logger {
	if WithHook && LoggerHook != nil {
		logger.LogEntry.Logger.AddHook(LoggerHook)
		return logger
	}
	return logger
}

func SetLogLevel(level string, logger *Logger) {
	switch level {
	case "debug":
		logger.SetLevel(log.DebugLevel)
	case "info":
		logger.SetLevel(log.InfoLevel)
	case "warn":
		logger.SetLevel(log.WarnLevel)
	case "error":
		logger.SetLevel(log.ErrorLevel)
	case "fatal":
		logger.SetLevel(log.FatalLevel)
	case "panic":
		logger.SetLevel(log.PanicLevel)
	default:
		logger.SetLevel(log.InfoLevel)
	}
}

func SetOutput(output string, logger *Logger) {
	switch output {
	case "stderr":
		logger.SetOutput(os.Stderr)
	case "stdout":
		logger.SetOutput(os.Stdout)
	case "null":
		logger.SetOutput(ioutil.Discard)
	default:
		logger.SetOutput(os.Stderr)
	}
}

func SetFormat(format string, logger *Logger) {
	switch format {
	case "json":
		logger.LogEntry.Logger.Formatter = &log.JSONFormatter{
			DisableTimestamp: true,
		}
	case "text":
		logger.LogEntry.Logger.Formatter = &log.TextFormatter{
			DisableTimestamp: true,
		}
	default:
		logger.LogEntry.Logger.Formatter = &log.JSONFormatter{
			DisableTimestamp: true,
		}
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	initLogger(logger).Debug(args...)
}

func (logger *Logger) Info(args ...interface{}) {
	initLogger(logger).Info(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	initLogger(logger).Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	initLogger(logger).Error(args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	initLogger(logger).Fatal(args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	initLogger(logger).Panic(args...)
}

// Entry Printf family functions
func (logger *Logger) Debugf(format string, args ...interface{}) {
	initLogger(logger).Debugf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	initLogger(logger).Infof(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	initLogger(logger).Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	initLogger(logger).Errorf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	initLogger(logger).Fatalf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	initLogger(logger).Panicf(format, args...)
}

func (logger *Logger) ReqId() string {
	reqId, ok := logger.LogEntry.Data[XReqIDConst]
	if ok {
		return fmt.Sprintf("%v", reqId)
	}
	return ""
}

func initLogger(logger *Logger) *log.Entry {
	entry := logger.LogEntry
	f := logger.fields
	return DecorateLog(entry.WithFields(log.Fields(f)).WithFields(log.Fields{"timedate": time.Now()}))
}

func (log *Logger) Xput(logs []string) {
	if xLog, exists := log.LogEntry.Data["X-Log"]; exists {
		if ll, ok := xLog.([]string); ok {
			ll = append(ll, logs...)
		}
	} else {
		log.LogEntry.Data["X-Log"] = logs
	}

}

func (log *Logger) ReqID() string {
	if reqid := log.ReqId(); len(reqid) > 0 {
		return reqid
	}

	return log.reqId
}
