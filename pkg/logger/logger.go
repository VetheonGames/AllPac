package logger

import (
	"io"
    "log"
    "os"
    "path/filepath"
	"fmt"
)

var Logger *log.Logger

func Init(logFilePath string) error {
    if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
        return err
    }

    logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }

    Logger = log.New(logFile, "AllPac: ", log.Ldate|log.Ltime|log.Lshortfile)
    Logger.SetOutput(io.MultiWriter(os.Stderr, logFile))

    return nil
}

func Info(v ...interface{}) {
    Logger.Println("INFO: " + fmt.Sprint(v...))
}

func Infof(format string, v ...interface{}) {
    Logger.Printf("INFO: "+format, v...)
}

func Warn(v ...interface{}) {
    Logger.Println("WARN: " + fmt.Sprint(v...))
}

func Warnf(format string, v ...interface{}) {
    Logger.Printf("WARN: "+format, v...)
}

func Error(v ...interface{}) {
    Logger.Println("ERROR: " + fmt.Sprint(v...))
}

func Errorf(format string, v ...interface{}) {
    Logger.Printf("ERROR: "+format, v...)
}

func Debug(v ...interface{}) {
    Logger.Println("DEBUG: " + fmt.Sprint(v...))
}

func Debugf(format string, v ...interface{}) {
    Logger.Printf("DEBUG: "+format, v...)
}
