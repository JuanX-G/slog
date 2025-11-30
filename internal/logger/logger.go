package httplogging

import (
	"log"
	"os"
	"time"
	"fmt"
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

func(l Logger) LogsWatchdog(toConsole, toFile bool, filename string) error {
	if !toFile {
		for v := range l.logChan {
			log.Print(v.msg)
		}
	} else if toFile && !toConsole{
		f, err := os.OpenFile(filename, os.O_WRONLY | os.O_CREATE, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		for v := range l.logChan {
			logString := fmt.Sprint("Time: ", v.timestamp, "Message: ", v.msg, "\n")
			f.Write([]byte(logString))
		}
	}
	return nil
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
