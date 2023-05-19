package logger

import (
	"fmt"
	"os"
	"time"
)

// singleton
var (
	loggerInstance *Logger
)

const (
	DEBUG = 1 << iota
	INFO
	WARN
	ERROR
	FATAL
)

type LogLevelType int
type Logger struct {
	logFileName string
	logLevel    LogLevelType
	logFile     *os.File
	err         error
}

func Init(logFileName string, logLevel LogLevelType) {
	var logger = Logger{}
	logger.logLevel = logLevel
	logger.logFileName = logFileName
	var logFile, err = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Can not open log file: " + logFileName)
		panic(err)
	}
	logger.logFile = logFile
	loggerInstance = &logger
}

func SetLogfileName(logFileName string) {
	loggerInstance.logFileName = logFileName
}

func SetLogLevel(logLevel LogLevelType) {
	loggerInstance.logLevel = logLevel
}

func Fatal(message string) {
	if loggerInstance.logLevel <= FATAL {
		loggerInstance.writeMessage("FATAL", message)
	}
}

func Error(message string) {
	if loggerInstance.logLevel <= ERROR {
		loggerInstance.writeMessage("ERROR", message)
	}
}

func Warn(message string) {
	if loggerInstance.logLevel <= WARN {
		loggerInstance.writeMessage("WARN", message)
	}
}

func Info(message string) {
	if loggerInstance.logLevel <= INFO {
		loggerInstance.writeMessage("INFO", message)
	}
}

func Debug(message string) {
	if loggerInstance.logLevel <= DEBUG {
		loggerInstance.writeMessage("DEBUG", message)
	}
}

func (logger *Logger) writeMessage(logLevel string, message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	loggerInstance.logFile.WriteString(currentTime + " " + logLevel + ": " + message + "\n")
}
