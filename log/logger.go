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
	LevelInfo  _level = 0
	LevelError _level = 1

	MemoryMaxLogLine = 128
	ErrorLogFile     = "error.log"
	InfoLogFile      = "info.log"
)

func isLogLevelValid(l _level) bool {
	if l != LevelInfo && l != LevelError {
		return false
	}

	return true
}

type Logger struct {
	logLevel _level

	infoLog   chan string
	infoText  string
	infoLines int

	errorLog   chan string
	errorText  string
	errorLines int

	storePath string // "" -> means no store log file
}

var logger *Logger

type Option func(lg *Logger)

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
		dir, _ := filepath.Split(storePath)
		if _, err := os.Stat(dir); err != nil {
			os.MkdirAll(dir, 0644)
		}
		lg.storePath = dir
	}
}

func NewLogger(opts ...Option) *Logger {
	if logger == nil {
		logger = &Logger{
			infoLog:  make(chan string, 256),
			errorLog: make(chan string, 256),
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
	for {
		select {
		case errText := <-l.errorLog:
			l.errorLines++
			if l.errorLines > MemoryMaxLogLine {
				l.errorText = "" // clear
			}
			l.errorText += errText

			if l.storePath != "" {
				errFile, err := os.OpenFile(filepath.Join(l.storePath, ErrorLogFile), os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					continue
				}
				errFile.WriteString(errText)
				errFile.Close()
			}
		case infoText := <-l.infoLog:
			if l.logLevel >= LevelError {
				continue
			}

			l.infoLines++
			if l.infoLines > MemoryMaxLogLine {
				l.infoText = ""
			}
			l.infoText += infoText

			if l.storePath != "" {
				errFile, err := os.OpenFile(filepath.Join(l.storePath, InfoLogFile), os.O_APPEND|os.O_CREATE, 0644)
				if err != nil {
					continue
				}
				errFile.WriteString(infoText)
				errFile.Close()
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

func InfoLogChan() <-chan string {
	return GetLogger().infoLog
}

func ErrorLogChan() <-chan string {
	return GetLogger().errorLog
}
