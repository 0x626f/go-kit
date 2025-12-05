package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

type Logger struct {
	name             string
	out, err         io.Writer
	syncOut, syncErr sync.Mutex
	options          *loggerOptions
}

type jsonLog struct {
	Source    string `json:"source,omitempty"`
	Level     string `json:"level,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Message   string `json:"message,omitempty"`
	Object    any    `json:"object,omitempty"`
}

func NewLogger() *Logger {
	return &Logger{
		out: os.Stdout,
		err: os.Stderr,
		options: &loggerOptions{
			level:           globalLoggerOptions.level,
			coloring:        globalLoggerOptions.coloring,
			timestamp:       globalLoggerOptions.timestamp,
			timestampFormat: globalLoggerOptions.timestampFormat,
		},
	}
}

func (logger *Logger) formatMessage(level LogLevel, msg string, args ...any) []byte {
	message := logger.format(msg, args...)
	if level == NONE {
		return append(message, ln)
	}

	var payload []byte

	if logger.options.timestamp {
		timestamp := time.Now().Format(logger.options.timestampFormat)
		payload = append(payload, []byte(timestamp)...)
		payload = append(payload, space)
	}

	payload = append(payload, []byte(level.String())...)
	payload = append(payload, space)

	if len(logger.name) > 0 {
		payload = append(payload, []byte(logger.name)...)
		payload = append(payload, space)
	}

	payload = append(payload, message...)

	if logger.options.coloring {
		payload = level.paint(payload)
	}

	return append(payload, ln)
}

func (logger *Logger) formatJSONMessage(level LogLevel, object any, msg string, args ...any) ([]byte, error) {
	log := jsonLog{
		Level:   level.String(),
		Message: fmt.Sprintf(msg, args...),
		Object:  object,
	}

	if level == NONE {
		raw, err := json.Marshal(log)
		if err != nil {
			return nil, err
		}
		return append(raw, ln), nil
	}

	if logger.options.timestamp {
		log.Timestamp = time.Now().Format(logger.options.timestampFormat)
	}

	if len(logger.name) > 0 {
		log.Source = logger.name
	}

	raw, err := json.Marshal(log)
	if err != nil {
		return nil, err
	}

	return append(raw, ln), nil
}

func (logger *Logger) format(msg string, args ...any) []byte {
	if len(args) == 0 {
		return []byte(msg)
	}
	return []byte(fmt.Sprintf(msg, args...))
}

func (logger *Logger) writeStringToStream(stream io.Writer, level LogLevel, msg string, args ...any) {
	if stream == nil {
		return
	}
	_, _ = stream.Write(logger.formatMessage(level, msg, args...))
}

func (logger *Logger) writeJSONToStream(stream io.Writer, level LogLevel, object any, msg string, args ...any) error {
	if stream == nil {
		return errors.New("nil stream or object")
	}

	payload, err := logger.formatJSONMessage(level, object, msg, args...)

	if err != nil {
		return err
	}

	_, err = stream.Write(payload)

	return err
}

func (logger *Logger) writeByLevel(level LogLevel, data []byte) {
	if level == ERROR {
		logger.syncErr.Lock()
		defer logger.syncErr.Unlock()
		_, _ = logger.err.Write(data)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()
	_, _ = logger.out.Write(data)
}

func (logger *Logger) sendToChannelByLevel(level LogLevel, data []byte) {
	if level == ERROR {
		logger.options.errors <- data
		return
	}

	logger.options.logs <- data
}

func (logger *Logger) OutputTo(target io.Writer) *Logger {
	if target != nil {
		logger.out = target
	}
	return logger
}

func (logger *Logger) ErrorsTo(target io.Writer) *Logger {
	if target != nil {
		logger.err = target
	}
	return logger
}

func (logger *Logger) WithAsync(option bool, capacity int) (*Logger, func()) {
	cancel := func() {}
	if option {
		logger.options.async = true
		logger.options.logs, logger.options.errors = make(chan []byte, capacity), make(chan []byte, capacity)
		logger.options.cancelAsync = make(chan struct{})

		cancel = func() {
			close(logger.options.cancelAsync)
		}

		go func(logs, errs chan []byte, cancel chan struct{}) {
			for {
				select {
				case log := <-logs:
					logger.syncOut.Lock()
					_, _ = logger.out.Write(log)
					logger.syncOut.Unlock()
				case err := <-errs:
					logger.syncErr.Lock()
					_, _ = logger.err.Write(err)
					logger.syncErr.Unlock()
				case <-cancel:
					return
				}
			}
		}(logger.options.logs, logger.options.errors, logger.options.cancelAsync)

	}
	return logger, cancel
}

func (logger *Logger) WithName(name string) *Logger {
	logger.name = name
	return logger
}

func (logger *Logger) WithLogLevel(level LogLevel) *Logger {
	logger.options.level = level
	return logger
}

func (logger *Logger) WithTimestamp() *Logger {
	logger.options.timestamp = true
	logger.options.timestampFormat = "2006-01-02 15:04:05"
	return logger
}

func (logger *Logger) WithTimestampFormat(format string) *Logger {
	logger.options.timestamp = true
	logger.options.timestampFormat = format
	return logger
}

func (logger *Logger) WithColoring() *Logger {
	if runtime.GOOS != "windows" {
		logger.options.coloring = true
	}
	return logger
}

func (logger *Logger) Logf(msg string, args ...any) {
	if logger.options.async {
		logger.options.logs <- logger.formatMessage(NONE, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, NONE, msg, args...)
}

func (logger *Logger) Tracef(msg string, args ...any) {
	if logger.options.level < TRACE {
		return
	}

	if logger.options.async {
		logger.options.logs <- logger.formatMessage(TRACE, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, TRACE, msg, args...)
}

func (logger *Logger) Debugf(msg string, args ...any) {
	if logger.options.level < DEBUG {
		return
	}

	if logger.options.async {
		logger.options.logs <- logger.formatMessage(DEBUG, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, DEBUG, msg, args...)
}

func (logger *Logger) Infof(msg string, args ...any) {
	if logger.options.level < INFO {
		return
	}

	if logger.options.async {
		logger.options.logs <- logger.formatMessage(INFO, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, INFO, msg, args...)
}

func (logger *Logger) Warningf(msg string, args ...any) {
	if logger.options.level < WARNING {
		return
	}

	if logger.options.async {
		logger.options.logs <- logger.formatMessage(WARNING, msg, args...)
		return
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	logger.writeStringToStream(logger.out, WARNING, msg, args...)
}

func (logger *Logger) Errorf(msg string, args ...any) {
	if logger.options.async {
		logger.options.errors <- logger.formatMessage(ERROR, msg, args...)
		return
	}

	logger.syncErr.Lock()
	defer logger.syncErr.Unlock()

	logger.writeStringToStream(logger.err, ERROR, msg, args...)
}

func (logger *Logger) LogJSONf(object any, msg string, args ...any) error {
	if logger.options.async {
		data, err := logger.formatJSONMessage(NONE, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, NONE, object, msg, args...)
}

func (logger *Logger) TraceJSONf(object any, msg string, args ...any) error {
	if logger.options.level < TRACE {
		return nil
	}

	if logger.options.async {
		data, err := logger.formatJSONMessage(TRACE, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, TRACE, object, msg, args...)
}

func (logger *Logger) DebugJSONf(object any, msg string, args ...any) error {
	if logger.options.level < DEBUG {
		return nil
	}

	if logger.options.async {
		data, err := logger.formatJSONMessage(DEBUG, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, DEBUG, object, msg, args...)
}

func (logger *Logger) InfoJSONf(object any, msg string, args ...any) error {
	if logger.options.level < INFO {
		return nil
	}

	if logger.options.async {
		data, err := logger.formatJSONMessage(INFO, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, INFO, object, msg, args...)
}

func (logger *Logger) WarningJSONf(object any, msg string, args ...any) error {
	if logger.options.level < WARNING {
		return nil
	}

	if logger.options.async {
		data, err := logger.formatJSONMessage(WARNING, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.logs <- data
		return nil
	}

	logger.syncOut.Lock()
	defer logger.syncOut.Unlock()

	return logger.writeJSONToStream(logger.out, WARNING, object, msg, args...)
}

func (logger *Logger) ErrorJSONf(object any, msg string, args ...any) error {
	if logger.options.async {
		data, err := logger.formatJSONMessage(ERROR, object, msg, args...)
		if err != nil {
			return err
		}
		logger.options.errors <- data
		return nil
	}

	logger.syncErr.Lock()
	defer logger.syncErr.Unlock()

	return logger.writeJSONToStream(logger.err, ERROR, object, msg, args...)
}

func (logger *Logger) LogObjectf(msg string, args ...any) *ObjectLogBuilder {
	return newObjectLogBuilder(logger, NONE, logger.format(msg, args...))
}

func (logger *Logger) TraceObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.level < TRACE {
		return nil
	}

	return newObjectLogBuilder(logger, TRACE, logger.format(msg, args...))

}

func (logger *Logger) DebugObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.level < DEBUG {
		return nil
	}

	return newObjectLogBuilder(logger, DEBUG, logger.format(msg, args...))

}

func (logger *Logger) InfoObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.level < INFO {
		return nil
	}

	return newObjectLogBuilder(logger, INFO, logger.format(msg, args...))

}

func (logger *Logger) WarningObjectf(msg string, args ...any) *ObjectLogBuilder {
	if logger.options.level < WARNING {
		return nil
	}

	return newObjectLogBuilder(logger, WARNING, logger.format(msg, args...))

}

func (logger *Logger) ErrorObjectf(msg string, args ...any) *ObjectLogBuilder {
	return newObjectLogBuilder(logger, ERROR, logger.format(msg, args...))

}
