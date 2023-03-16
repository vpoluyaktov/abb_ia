package logger

import (
	"os"
  "time"
)

var (
  logLevelsOrder []string = []string{"FATAL", "ERROR", "WARN", "INFO", "DEBUG"}
  currentlogFileName string
  currentlogLevel string
  logFile *os.File
  err error
)

func Init(logFileName string, logLevel string) {
  currentlogLevel = logLevel
  currentlogFileName = logFileName
  logFile, err = os.Create(logFileName)
  if err != nil {
    panic(err)
  }
}

func Fatal(message string) {
	writeMessage("FATAL", message)
}

func Error(message string) {
	writeMessage("ERROR", message)
}

func Warn(message string) {
	writeMessage("WARN", message)
}

func Info(message string) {
	writeMessage("INFO", message)
}

func Debug(message string) {
	writeMessage("DEBUG", message)
}

func writeMessage(logLevel string, message string) {
  for _, level := range logLevelsOrder {
    if level == logLevel {
      currentTime := time.Now().Format("2006-01-02 15:04:05")
      logFile.WriteString(currentTime + " " + logLevel + ": " + message + "\n")       
    }
    if level == currentlogLevel {
      break
    }
  }
}