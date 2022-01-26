package log

import (
	"github.com/op/go-logging"
	"strings"
)

func Highlight(msg string, level logging.Level) string {
	buf := strings.Builder{}
	switch level {
	case logging.ERROR:
		buf.WriteString("<div class='log-error'>")
		buf.WriteString(msg)
		buf.WriteString("</div>")
		return buf.String()
	case logging.DEBUG:
		buf.WriteString("<div class='log-debug'>")
		buf.WriteString(msg)
		buf.WriteString("</div>")
		return buf.String()
	case logging.WARNING:
		buf.WriteString("<div class='log-warning'>")
		buf.WriteString(msg)
		buf.WriteString("</div>")
		return buf.String()
	default:
		buf.WriteString("<div class='log-info'>")
		buf.WriteString(msg)
		buf.WriteString("</div>")
		return buf.String()
	}
}
