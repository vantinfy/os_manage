package log

import (
	"fmt"
	"os"
	"os_manage/config"
	"path/filepath"
	"syscall"
	"time"
)

type _level int

const (
	LevelDebug _level = 0
	LevelInfo  _level = 1
	LevelError _level = 3

	ErrorLogFile = "error.log"
	InfoLogFile  = "info.log"
	DebugLogFile = "debug.log"
)

func isLogLevelValid(l _level) bool {
	if l != LevelDebug && l != LevelInfo && l != LevelError {
		return false
	}

	return true
}

type Logger struct {
	logLevel _level

	debugLog chan string
	infoLog  chan string
	errorLog chan string

	storePath string // "" -> means no store log file

	Extend        any
	ExtendHandler func(any, string)
}

var logger *Logger

type Option func(lg *Logger)

func WithLogExtend(extend any, extendHandler func(eha any, msg string)) Option {
	return func(lg *Logger) {
		lg.Extend = extend
		lg.ExtendHandler = extendHandler
	}
}

func WithLogLevel(logLevel int) Option {
	return func(lg *Logger) {
		ll := _level(logLevel)
		if !isLogLevelValid(ll) {
			ll = LevelInfo
		}
		lg.logLevel = ll
	}
}

func WithStorePath(storePath string) Option {
	return func(lg *Logger) {
		if _, err := os.Stat(storePath); err != nil {
			_ = os.MkdirAll(storePath, 0644)
		}
		lg.storePath = storePath
	}
}

func NewLogger(opts ...Option) *Logger {
	if logger == nil {
		logger = &Logger{
			debugLog: make(chan string, 10),
			infoLog:  make(chan string, 10),
			errorLog: make(chan string, 10),
		}
	}
	for _, opt := range opts {
		opt(logger)
	}

	go logger.handleLogChan()

	return logger
}

func GetLogger() *Logger {
	if logger == nil {
		logger = NewLogger()
	}

	return logger
}

func (l *Logger) handleLogChan() {
	var (
		errFile   *os.File
		infoFile  *os.File
		debugFile *os.File
		err       error
	)
	// 注意参数f是二级指针
	openLogFile := func(name string, f **os.File) {
		*f, err = os.OpenFile(filepath.Join(l.storePath, name), os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			panic(fmt.Sprintf("open log file[%s] failed: %v", filepath.Join(l.storePath, name), err))
		}
	}
	if l.storePath != "" {
		openLogFile(ErrorLogFile, &errFile)
		defer errFile.Close()

		if l.logLevel <= LevelInfo {
			openLogFile(InfoLogFile, &infoFile)
			defer infoFile.Close()
		}

		if l.logLevel <= LevelDebug {
			openLogFile(DebugLogFile, &debugFile)
			defer debugFile.Close()
		}
	}

	for {
		select {
		case errText := <-l.errorLog:
			if l.ExtendHandler != nil {
				l.ExtendHandler(l.Extend, errText)
			}

			if errFile != nil {
				errFile.WriteString(errText)
			}

		case infoText := <-l.infoLog:
			if l.ExtendHandler != nil {
				l.ExtendHandler(l.Extend, infoText)
			}

			if l.logLevel > LevelInfo {
				continue
			}

			if infoFile != nil {
				infoFile.WriteString(infoText)
			}

		case debugText := <-l.debugLog:
			if l.ExtendHandler != nil {
				l.ExtendHandler(l.Extend, debugText)
			}

			if l.logLevel > LevelDebug {
				continue
			}

			if debugFile != nil {
				debugFile.WriteString(debugText)
			}
		}
	}
}

func (l *Logger) Fatal(logContent ...any) {
	appendContent := fmt.Sprintf("[Fatal] %s: %v\n",
		time.Now().Format("2006-01-02 15:04:05"), logContent)

	l.errorLog <- appendContent

	time.AfterFunc(time.Millisecond*3, func() {
		config.GlobalQuit <- syscall.SIGQUIT
	})
}

func (l *Logger) Error(logContent ...any) {
	appendContent := fmt.Sprintf("[Error] %s: %v\n",
		time.Now().Format("2006-01-02 15:04:05"), logContent)

	l.errorLog <- appendContent
}

func (l *Logger) Info(logContent ...any) {
	appendContent := fmt.Sprintf("[Info] %s: %v\n",
		time.Now().Format("2006-01-02 15:04:05"), logContent)

	l.infoLog <- appendContent
}

func (l *Logger) Debug(logContent ...any) {
	appendContent := fmt.Sprintf("[Debug] %s: %v\n",
		time.Now().Format("2006-01-02 15:04:05"), logContent)

	l.debugLog <- appendContent
}

func Fatal(logContent ...any) {
	GetLogger().Fatal(logContent)
}

func Error(logContent ...any) {
	GetLogger().Error(logContent)
}

func Errorf(format string, logContent ...any) {
	GetLogger().Error(fmt.Sprintf(format, logContent))
}

func Info(logContent ...any) {
	GetLogger().Info(logContent)
}

func Debug(logContent ...any) {
	GetLogger().Debug(logContent)
}

func DebugLogChan() <-chan string {
	return GetLogger().debugLog
}

func InfoLogChan() <-chan string {
	return GetLogger().infoLog
}

func ErrorLogChan() <-chan string {
	return GetLogger().errorLog
}
