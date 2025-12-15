package msg

import (
	"fmt"
	"strings"
	"time"
	"os"
)

const (
	LvlFatal = iota + 1
	LvlError
	LvlWarning
	LvlInfo
	LvlVerbose
	LvlVerbose2
)

var DoTimestamp bool = os.Getenv("JOURNAL_STREAM") == ""
var Level int = LvlInfo
var startTime time.Time

func init() {
	startTime = time.Now()
}

func domsg(level int, format string, a ...any) {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("%d ", level))
	if DoTimestamp {
		duration := time.Since(startTime)
		s.WriteRune('[')
		s.WriteString(fmt.Sprintf("%6d.%06d", duration.Microseconds()/1e6, duration.Microseconds()%1e6))
		s.WriteRune(']')
		s.WriteRune(' ')
	}
	s.WriteString(fmt.Sprintf(format, a...))
	if level == LvlFatal {
		panic(s.String())
	}
	if Level >= level {
		fmt.Println(s.String())
	}
}

func Fatal(format string, a ...any) {
	domsg(LvlFatal, format, a...)
}
func Error(format string, a ...any) {
	domsg(LvlError, format, a...)
}
func Warning(format string, a ...any) {
	domsg(LvlWarning, format, a...)
}
func Info(format string, a ...any) {
	domsg(LvlInfo, format, a...)
}
func Verbose(format string, a ...any) {
	domsg(LvlVerbose, format, a...)
}
func Verbose2(format string, a ...any) {
	domsg(LvlVerbose2, format, a...)
}
