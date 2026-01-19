package stdoutWriter
import (
	"bufio"
	"log"
	"os"
)
// -----------------------------
// Console writer
// -----------------------------
type StdoutWriter struct {
	w     *bufio.Writer
	queue chan string
}

func NewStdoutWriter() *StdoutWriter {
	BWriter := bufio.NewWriter(os.Stdout)
	sw := &StdoutWriter{
		w:     BWriter,
		queue: make(chan string, 128),
	}
	go func(writer *StdoutWriter) {
		for msg := range writer.queue {
			_, err := writer.w.WriteString(msg + "\n")
			if err != nil {
				log.Panic("write error: ", err)
			}
			writer.w.Flush()
		}
	}(sw)
	return sw
}

func (w *StdoutWriter) AppendQueue(s ...string) {
	msg := ""
	for _, v := range s {
		msg += v
	}
	w.queue <- msg
}
