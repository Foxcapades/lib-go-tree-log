package tlog

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// TreeLogger Indentation Defaults
const (
	indDefIndent = "  "
	indDefInLen  = uint8(len(indDefIndent))
	indDefCache  = "                    "
	indDefChLen  = uint8(len(indDefCache)) / indDefInLen
)

// TreeLogger provides functions to print logs in a
// hierarchical presentation to help provide clarity for
// verbose processes.
type TreeLogger interface {

	IndentString(txt string) TreeLogger

	// Increases the indentation for the logger.
	//
	// Each call of this will increase the indentation by 1
	// copy of the IndentString value for following write
	// calls.
	Indent() TreeLogger

	// Decreases the indentation for the logger.
	//
	// Each call of this will decrease the indentation by 1
	// copy of the IndentString value for following write
	// calls.
	UnIndent() TreeLogger

	// Sets the writer into which logs will be written.
	//
	// Default writer is os.Stdout.
	Writer(out io.Writer) TreeLogger

	// Writes the given values at the current indentation
	// level to the log writer with no trailing newline.
	Write(in ...interface{}) TreeLogger

	// Writes the given values as a single line at the current
	// indentation level to the log writer.
	WriteLn(in ...interface{}) TreeLogger

	// Writes the given values as a single line at an
	// increased indentation level to the log writer.
	//
	// Shortcut for calling Indent().WriteLn(...).UnIndent()
	WriteChildLn(in ...interface{}) TreeLogger

	// Directly appends the given values to the current log
	// line without prepending the current indentation or
	// appending a trailing newline
	Append(in ...interface{}) TreeLogger

	// Prints a newline character to the log writer
	NewLine() TreeLogger
}

// NewTreeLogger creates a new default configured
// implementation of TreeLogger.
func NewTreeLogger() TreeLogger {
	return &treeLogger{
		indent:  indDefIndent,
		inLen:   indDefInLen,
		cache:   indDefCache,
		chLen:   indDefChLen,
		level:   0,
		writer:  os.Stdout,
		current: []interface{}{""},
	}
}

var defLogger = NewTreeLogger()

// DefaultLogger returns the preconfigured default logger.
//
// All calls to this method will return the same TreeLogger
// instance.
func DefaultLogger() TreeLogger {
	return defLogger
}

type treeLogger struct {
	cache string
	indent string

	current []interface{}

	writer io.Writer

	level uint8
	chLen uint8
	inLen uint8
}

func (t *treeLogger) Writer(out io.Writer) TreeLogger {
	t.writer = out
	return t
}

func (t *treeLogger) IndentString(txt string) TreeLogger {
	t.indent = txt
	return t
}

func (t *treeLogger) Indent() TreeLogger {
	t.level++
	t.clipIndent()
	return t
}

func (t *treeLogger) UnIndent() TreeLogger {
	t.level--
	t.clipIndent()
	return t
}

func (t *treeLogger) Write(in ...interface{}) TreeLogger {
	forceWrite(t.writer, t.current)
	forceWrite(t.writer, in)
	return t
}

func (t *treeLogger) WriteLn(in ...interface{}) TreeLogger {
	forceWrite(t.writer, t.current)
	forceWriteLn(t.writer, in)
	return t
}

func (t *treeLogger) WriteChild(in ...interface{}) TreeLogger {
	t.Indent()
	t.Write(in)
	t.UnIndent()
	return t
}

func (t *treeLogger) WriteChildLn(in ...interface{}) TreeLogger {
	t.Indent()
	t.WriteLn(in)
	t.UnIndent()
	return t
}

func (t *treeLogger) Append(in ...interface{}) TreeLogger {
	forceWrite(t.writer, in)
	return t
}

func (t *treeLogger) NewLine() TreeLogger {
	forceWriteLn(t.writer, nil)
	return t
}

func (t *treeLogger) resizeCache() {
	t.chLen = t.level
	t.cache = strings.Repeat(t.indent, int(t.level))
}

func (t *treeLogger) clipIndent() {
	if t.level > t.chLen {
		t.resizeCache()
	}
	t.current[0] = t.cache[:t.inLen*t.level]
}

func forceWrite(w io.Writer, in []interface{}) {
	forcePush(w, in, fmt.Fprint)
}

func forceWriteLn(w io.Writer, in []interface{}) {
	forcePush(w, in, fmt.Fprintln)
}

type pusher func(io.Writer, ...interface{}) (int, error)

func forcePush(w io.Writer, in []interface{}, fn pusher) {
	var e error
	if in != nil {
		_, e = fn(w, in...)
	} else {
		_, e = fn(w)
	}

	if e != nil {
		panic(e)
	}
}
