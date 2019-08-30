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

type TreeLogger interface {
	IndentString(txt string) TreeLogger

	Indent() TreeLogger
	UnIndent() TreeLogger

	Writer(out io.Writer) TreeLogger

	Write(in ...interface{}) TreeLogger
	WriteLn(in ...interface{}) TreeLogger

	WriteChild(in ...interface{}) TreeLogger
	WriteChildLn(in ...interface{}) TreeLogger

	Append(in ...interface{}) TreeLogger
	NewLine() TreeLogger
}

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

func DefaultLogger() TreeLogger {
	return defLogger
}

type treeLogger struct {
	cache string
	chLen uint8

	current []interface{}

	indent string
	inLen  uint8

	level uint8

	writer io.Writer
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
