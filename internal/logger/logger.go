package httplogging

import (
	"os"
	"time"
	"fmt"

	"slog-simple-blog/internal/stdoutWriter"
)

type logEntry struct {
	timestamp time.Time 
	msg string
}

type Logger struct {
	logChan chan logEntry
	file *os.File
	writer *stdoutWriter.StdoutWriter
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
			logString := fmt.Sprint("Time: ", v.timestamp, "Message: ", v.msg, "\n")
			l.writer.AppendQueue(logString)
		}
	} else if toFile && !toConsole && l.file != nil {
		for v := range l.logChan {
			logString := fmt.Sprint("Time: ", v.timestamp, "Message: ", v.msg, "\n")
			l.file.Write([]byte(logString))
		}
	} else if toFile && toConsole && l.file != nil {
		for v := range l.logChan {
			logString := fmt.Sprint("Time: ", v.timestamp, "Message: ", v.msg, "\n")
			l.file.Write([]byte(logString))
			l.writer.AppendQueue(logString)
		}
	}
	return nil
}

func (l Logger) Close() {
	close(l.logChan)
}

func NewLogger(toFile bool, fileName string) (*Logger, error) {
	if !toFile {
		return &Logger{
			logChan: make(chan logEntry, 100),
			file: nil,
			writer: stdoutWriter.NewStdoutWriter(),
		}, nil
	} else {
		file, err := os.Open(fileName)
		if err != nil {return nil, err}
		return &Logger{
			logChan: make(chan logEntry, 100),
			file: file,
			writer: stdoutWriter.NewStdoutWriter(),
		}, nil
	}
}

type Logable interface {
	String() string
}
