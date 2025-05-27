package log

import (
	"fmt"
	"io"
)

type Log struct {
	Out     io.Writer
	verbose bool
}

func NewLog(out io.Writer, verbose bool) *Log {
	return &Log{
		Out:     out,
		verbose: verbose,
	}
}

func (l *Log) Log(msg string) {
	if l.verbose {
		fmt.Fprintf(l.Out, "%s", msg)
	}
}
