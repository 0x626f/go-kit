package console

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Console struct {
	name                     string
	in                       io.Reader
	out, err                 io.Writer
	syncIn, syncOut, syncErr sync.Mutex
	options                  *consoleOptions
}

func NewConsole() *Console {
	return &Console{
		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
		options: &consoleOptions{
			level:           globalConsoleOptions.level,
			coloring:        globalConsoleOptions.coloring,
			timestamp:       globalConsoleOptions.timestamp,
			timestampFormat: globalConsoleOptions.timestampFormat,
		},
	}
}

func (console *Console) formatMessage(level LogLevel, msg string, args ...any) []byte {
	message := console.format(msg, args...)
	if level == NONE {
		return append(message, ln)
	}

	var payload []byte

	if console.options.timestamp {
		timestamp := time.Now().Format(console.options.timestampFormat)
		payload = append(payload, []byte(timestamp)...)
		payload = append(payload, space)
	}

	payload = append(payload, []byte(level.String())...)
	payload = append(payload, space)

	if len(console.name) > 0 {
		payload = append(payload, []byte(console.name)...)
		payload = append(payload, space)
	}

	payload = append(payload, message...)

	if console.options.coloring {
		payload = level.paint(payload)
	}

	return append(payload, ln)
}

func (console *Console) format(msg string, args ...any) []byte {
	if len(args) == 0 {
		return []byte(msg)
	}
	return []byte(fmt.Sprintf(msg, args...))
}

func (console *Console) writeStringToStream(stream io.Writer, level LogLevel, msg string, args ...any) {
	if stream == nil {
		return
	}
	_, _ = stream.Write(console.formatMessage(level, msg, args...))
}

func (console *Console) InputTo(target io.Reader) *Console {
	if target != nil {
		console.in = target
	}
	return console
}

func (console *Console) OutputTo(target io.Writer) *Console {
	if target != nil {
		console.out = target
	}
	return console
}

func (console *Console) ErrorsTo(target io.Writer) *Console {
	if target != nil {
		console.err = target
	}
	return console
}

func (console *Console) WithName(name string) *Console {
	console.name = name
	return console
}

func (console *Console) WithLogLevel(level LogLevel) *Console {
	console.options.level = level
	return console
}

func (console *Console) WithTimestamp() *Console {
	console.options.timestamp = true
	console.options.timestampFormat = "2006-01-02 15:04:05"
	return console
}

func (console *Console) WithTimestampFormat(format string) *Console {
	console.options.timestamp = true
	console.options.timestampFormat = format
	return console
}

func (console *Console) WithColoring() *Console {
	if runtime.GOOS != "windows" {
		console.options.coloring = true
	}
	return console
}

func (console *Console) Logf(msg string, args ...any) {
	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	console.writeStringToStream(console.out, NONE, msg, args...)
}

func (console *Console) Tracef(msg string, args ...any) {
	if console.options.level < TRACE {
		return
	}

	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	console.writeStringToStream(console.out, TRACE, msg, args...)
}

func (console *Console) Debugf(msg string, args ...any) {
	if console.options.level < DEBUG {
		return
	}

	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	console.writeStringToStream(console.out, DEBUG, msg, args...)
}

func (console *Console) Infof(msg string, args ...any) {
	if console.options.level < INFO {
		return
	}

	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	console.writeStringToStream(console.out, INFO, msg, args...)
}

func (console *Console) Warningf(msg string, args ...any) {
	if console.options.level < WARNING {
		return
	}

	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	console.writeStringToStream(console.out, WARNING, msg, args...)
}

func (console *Console) Errorf(msg string, args ...any) {
	console.syncErr.Lock()
	defer console.syncErr.Unlock()

	console.writeStringToStream(console.err, ERROR, msg, args...)
}

func (console *Console) Write(data []byte) (n int, err error) {
	if console.out == nil {
		return 0, errors.New("writer is nil")
	}

	console.syncOut.Lock()
	defer console.syncOut.Unlock()

	return console.out.Write(data)
}

func (console *Console) Read(buffer []byte) (n int, err error) {
	if console.in == nil {
		return 0, errors.New("reader is nil")
	}

	console.syncIn.Lock()
	defer console.syncIn.Unlock()

	return console.in.Read(buffer)
}

func (console *Console) ReadAll() ([]byte, error) {
	if console.in == nil {
		return nil, errors.New("reader is nil")
	}

	console.syncIn.Lock()
	defer console.syncIn.Unlock()

	return io.ReadAll(console.in)
}
