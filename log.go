// The colorful and simple logging library
// Copyright (c) 2017 Fadhli Dzil Ikram

package log

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/withmandala/go-log/colorful"
	"golang.org/x/crypto/ssh/terminal"
)

// FdWriter interface extends existing io.Writer with file descriptor function
// support
type FdWriter interface {
	io.Writer
	Fd() uintptr
}

// Logger struct define the underlying storage for single logger
type Logger struct {
	mu        sync.Mutex
	color     bool
	out       FdWriter
	debug     bool
	timestamp bool
	buf       colorful.ColorBuffer
}

// Prefix struct define plain and color byte
type Prefix struct {
	Plain []byte
	Color []byte
	File  bool
}

var (
	// Plain prefix template
	plainError = []byte("[ERROR] ")
	plainWarn  = []byte("[WARN]  ")
	plainInfo  = []byte("[INFO]  ")
	plainDebug = []byte("[DEBUG] ")
	plainTrace = []byte("[TRACE] ")

	// ErrorPrefix show error prefix
	ErrorPrefix = Prefix{
		Plain: plainError,
		Color: colorful.Red(plainError),
		File:  true,
	}

	// WarnPrefix show warn prefix
	WarnPrefix = Prefix{
		Plain: plainWarn,
		Color: colorful.Orange(plainWarn),
	}

	// InfoPrefix show info prefix
	InfoPrefix = Prefix{
		Plain: plainInfo,
		Color: colorful.Green(plainInfo),
	}

	// DebugPrefix show info prefix
	DebugPrefix = Prefix{
		Plain: plainDebug,
		Color: colorful.Cyan(plainDebug),
		File:  true,
	}

	// TracePrefix show info prefix
	TracePrefix = Prefix{
		Plain: plainTrace,
		Color: colorful.Purple(plainTrace),
	}
)

// New returns new Logger instance with predefined writer output and
// automatically detect terminal coloring support
func New(out FdWriter) *Logger {
	return &Logger{
		color:     terminal.IsTerminal(int(out.Fd())),
		out:       out,
		timestamp: true,
	}
}

// WithColor explicitly turn on colorful features on the log
func (l *Logger) WithColor() *Logger {
	l.mu.Lock()
	l.color = true
	l.mu.Unlock()
	return l
}

// WithoutColor explicitly turn off colorful features on the log
func (l *Logger) WithoutColor() *Logger {
	l.mu.Lock()
	l.color = false
	l.mu.Unlock()
	return l
}

// WithDebug turn on debugging output on the log to reveal debug and trace level
func (l *Logger) WithDebug() *Logger {
	l.mu.Lock()
	l.debug = true
	l.mu.Unlock()
	return l
}

// WithoutDebug turn off debugging output on the log
func (l *Logger) WithoutDebug() *Logger {
	l.mu.Lock()
	l.debug = false
	l.mu.Unlock()
	return l
}

// WithTimestamp turn on timestamp output on the log
func (l *Logger) WithTimestamp() *Logger {
	l.mu.Lock()
	l.timestamp = true
	l.mu.Unlock()
	return l
}

// WithoutTimestamp turn off timestamp output on the log
func (l *Logger) WithoutTimestamp() *Logger {
	l.mu.Lock()
	l.timestamp = false
	l.mu.Unlock()
	return l
}

// Output print the actual value
func (l *Logger) Output(depth int, prefix Prefix, data string) error {
	now := time.Now()
	var file string
	var line int
	if prefix.File {
		var ok bool
		if _, file, line, ok = runtime.Caller(depth + 1); !ok {
			file = "<unknown file>"
			line = 0
		}
	}
	// Acquire lock
	l.mu.Lock()
	defer l.mu.Unlock()
	// Reset buffer
	l.buf.Reset()
	// Add prefix
	if l.color {
		l.buf.Append(prefix.Color)
	} else {
		l.buf.Append(prefix.Plain)
	}
	// Add date-time
	if l.timestamp {
		if l.color {
			l.buf.Blue()
		}
		year, month, day := now.Date()
		l.cint(year, 4)
		l.buf.AppendByte('/')
		l.cint(int(month), 2)
		l.buf.AppendByte('/')
		l.cint(day, 2)
		l.buf.AppendByte(' ')
		hour, min, sec := now.Clock()
		l.cint(hour, 2)
		l.buf.AppendByte(':')
		l.cint(min, 2)
		l.buf.AppendByte(':')
		l.cint(sec, 2)
		l.buf.AppendByte(' ')
		if l.color {
			l.buf.Off()
		}
	}
	// Add file
	if prefix.File {
		if l.color {
			l.buf.Orange()
		}
		l.buf.Append([]byte(file))
		l.buf.AppendByte(':')
		l.cint(line, 0)
		l.buf.AppendByte(' ')
		if l.color {
			l.buf.Off()
		}
	}
	// Add data
	l.buf.Append([]byte(data))
	if len(data) == 0 || data[len(data)-1] != '\n' {
		l.buf.AppendByte('\n')
	}
	// Flush to output
	_, err := l.out.Write(l.buf.Buffer)
	return err
}

func (l *Logger) cint(val int, width int) {
	var repr [8]byte
	reprCount := len(repr) - 1
	for val >= 10 || width > 1 {
		reminder := val / 10
		repr[reprCount] = byte('0' + val - reminder*10)
		val = reminder
		reprCount--
		width--
	}
	repr[reprCount] = byte('0' + val)
	l.buf.Append(repr[reprCount:])
}

// Error print error message to output and quit the application with status 1
func (l *Logger) Error(v ...interface{}) {
	l.Output(1, ErrorPrefix, fmt.Sprintln(v...))
	os.Exit(1)
}

// Errorf print formatted error message to output and quit the application
// with status 1
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Output(1, ErrorPrefix, fmt.Sprintf(format, v...))
	os.Exit(1)
}

// Warn print warning message to output
func (l *Logger) Warn(v ...interface{}) {
	l.Output(1, WarnPrefix, fmt.Sprintln(v...))
}

// Warnf print formatted warning message to output
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Output(1, WarnPrefix, fmt.Sprintf(format, v...))
}

// Info print informational message to output
func (l *Logger) Info(v ...interface{}) {
	l.Output(1, InfoPrefix, fmt.Sprintln(v...))
}

// Infof print formatted informational message to output
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Output(1, InfoPrefix, fmt.Sprintf(format, v...))
}

// Debug print debug message to output if debug output enabled
func (l *Logger) Debug(v ...interface{}) {
	if l.debug {
		l.Output(1, DebugPrefix, fmt.Sprintln(v...))
	}
}

// Debugf print formatted debug message to output if debug output enabled
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.debug {
		l.Output(1, DebugPrefix, fmt.Sprintf(format, v...))
	}
}

// Trace print trace message to output if debug output enabled
func (l *Logger) Trace(v ...interface{}) {
	if l.debug {
		l.Output(1, TracePrefix, fmt.Sprintln(v...))
	}
}

// Tracef print formatted trace message to output if debug output enabled
func (l *Logger) Tracef(format string, v ...interface{}) {
	if l.debug {
		l.Output(1, TracePrefix, fmt.Sprintf(format, v...))
	}
}
