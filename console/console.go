package console

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

type LogLevel uint8

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
	TRACE
)

func (level LogLevel) String() string {
	switch level {
	case ERROR:
		return "ERROR"
	case WARNING:
		return "WARNING"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case TRACE:
		return "TRACE"
	}
	return ""
}

type consoleOptions struct {
	level LogLevel

	timestamp       bool
	timestampFormat string
}

type Console struct {
	name         string
	in, out, err *os.File
	options      *consoleOptions
}

func NewConsole() *Console {
	return &Console{
		in:      os.Stdin,
		out:     os.Stdout,
		err:     os.Stderr,
		options: &consoleOptions{level: TRACE},
	}
}

func (console *Console) formatMessage(level LogLevel, msg string) string {
	name := ""

	if len(console.name) != 0 {
		name = fmt.Sprintf(" - [ %v ]", console.name)
	}

	timestamp := ""
	if console.options.timestamp {
		now := time.Now()
		timestamp = fmt.Sprintf(" - [ %v ]", now.Format("01/02/2006 3:04:05 PM"))

		if len(console.options.timestampFormat) != 0 {
			timestamp = fmt.Sprintf(" - [ %v ]", now.Format(console.options.timestampFormat))
		}
	}
	return fmt.Sprintf("[ %v ]%v%v: %v\n", level.String(), timestamp, name, msg)
}

func (console *Console) InputTo(target *os.File) *Console {
	if target != nil {
		console.in = target
	}
	return console
}

func (console *Console) OutputTo(target *os.File) *Console {
	if target != nil {
		console.out = target
	}
	return console
}

func (console *Console) ErrorsTo(target *os.File) *Console {
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
	return console
}

func (console *Console) WithTimestampFormat(format string) *Console {
	console.options.timestamp = true
	console.options.timestampFormat = format
	return console
}

func (console *Console) Log(msg string) {
	_, _ = console.out.Write([]byte(msg + "\n"))
}

func (console *Console) Trace(msg string) {
	if console.options.level < TRACE {
		return
	}
	_, _ = console.out.Write([]byte(console.formatMessage(TRACE, msg)))
}

func (console *Console) Debug(msg string) {
	if console.options.level < DEBUG {
		return
	}
	_, _ = console.out.Write([]byte(console.formatMessage(DEBUG, msg)))
}

func (console *Console) Info(msg string) {
	if console.options.level < INFO {
		return
	}
	_, _ = console.out.Write([]byte(console.formatMessage(INFO, msg)))
}

func (console *Console) Warning(msg string) {
	if console.options.level < WARNING {
		return
	}
	_, _ = console.out.Write([]byte(console.formatMessage(WARNING, msg)))
}

func (console *Console) Error(msg string) {
	_, _ = console.out.Write([]byte(console.formatMessage(ERROR, msg)))
}

func (console *Console) Logf(msg string, args ...any) {
	console.Log(fmt.Sprintf(msg, args...))
}

func (console *Console) Tracef(msg string, args ...any) {
	console.Trace(fmt.Sprintf(msg, args...))
}

func (console *Console) Debugf(msg string, args ...any) {
	console.Debug(fmt.Sprintf(msg, args...))
}

func (console *Console) Infof(msg string, args ...any) {
	console.Info(fmt.Sprintf(msg, args...))
}

func (console *Console) Warningf(msg string, args ...any) {
	console.Warning(fmt.Sprintf(msg, args...))
}

func (console *Console) Errorf(msg string, args ...any) {
	console.Error(fmt.Sprintf(msg, args...))
}

func (console *Console) Read(n int) (int, []byte, error) {
	buffer := make([]byte, n)
	length, err := console.in.Read(buffer)
	return length, buffer, err
}

func (console *Console) ReadAll() ([]byte, error) {
	return io.ReadAll(console.in)
}

func (console *Console) Reader() *bufio.Reader {
	return bufio.NewReader(console.in)
}

type ObjectLog struct {
	console *Console
}
