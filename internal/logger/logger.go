package httplogging

import (
	"log"
	"time"
)

type logEntry struct {
	timestamp time.Time 
	msg string
}

type Logger struct {
	logChan chan logEntry
}

func (l *Logger) LogString(msg string) {
    l.logChan <- logEntry{msg: msg, timestamp: time.Now()} 
}

func (l *Logger) LogStruct(arg Logable) {
    l.logChan <- logEntry{msg: arg.String(), timestamp: time.Now()} 
}

func(l Logger) LogsWatchdog() {
	for v := range l.logChan {
		log.Print(v.msg)
	}
}

func (l Logger) Close() {
	close(l.logChan)
}

func NewLogger() *Logger {
	return &Logger{
		logChan: make(chan logEntry, 100),
	}
}


type Logable interface {
	String() string
}
