package core

import (
	"fmt"
	"log"
)

type Logger func(format string, v ...interface{})

var gLogger Logger

func SetLogger(l Logger) {
	gLogger = l
}
func GetLogger() Logger {
	return gLogger
}
func DoLog(format string, v ...interface{}) {
	if gLogger != nil {
		gLogger(format, v...)
	}
}

func LoggerGo(format string, v ...interface{}) {
	log.Printf(format+"\n", v...)
}

func LoggerFmtPrint(format string, v ...interface{}) {
	fmt.Printf(format+"\n", v...)
}
