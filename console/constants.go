package console

type LogLevel uint8
type color []byte

const (
	ERROR LogLevel = iota
	WARNING
	INFO
	DEBUG
	TRACE
	NONE

	space = byte(' ')
	ln    = byte('\n')
)

var (
	colorReset  color = []byte{27, 91, 48, 109}
	colorRed    color = []byte{27, 91, 51, 49, 109}
	colorGreen  color = []byte{27, 91, 51, 50, 109}
	colorYellow color = []byte{27, 91, 51, 51, 109}
	colorBlue   color = []byte{27, 91, 51, 52, 109}
	colorGrey   color = []byte{27, 91, 57, 48, 109}
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
	default:
		return ""
	}
}

func (level LogLevel) color() color {
	switch level {
	case ERROR:
		return colorRed
	case WARNING:
		return colorYellow
	case INFO:
		return colorGreen
	case DEBUG:
		return colorGrey
	case TRACE:
		return colorBlue
	default:
		return []byte{}
	}
}

func (level LogLevel) paint(payload []byte) []byte {
	return append(append(level.color(), payload...), colorReset...)
}
