package printing

import (
	"io"
	"runtime"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
)

type colors struct {
	out   io.Writer
	inner Printer
}

func NewColorPrinter(out io.Writer, inner Printer) *colors {
	return &colors{out: out, inner: inner}
}

func (this *colors) Print(token lexing.Token) {
	// TODO: string values: green
	// TODO: object keys: blue
	// TODO: array and object braces: cyan?
	switch token.Type {
	case lexing.TokenNull:
		this.write(Gray)
		this.inner.Print(token)
		this.write(Reset)
	default:
		this.inner.Print(token)
	}
}

func (this *colors) write(data []byte) {
	_, _ = this.out.Write(data)
}

var (
	Reset  = []byte("\033[0m")
	Red    = []byte("\033[31m")
	Green  = []byte("\033[32m")
	Yellow = []byte("\033[33m")
	Blue   = []byte("\033[34m")
	Purple = []byte("\033[35m")
	Cyan   = []byte("\033[36m")
	Gray   = []byte("\033[37m")
	White  = []byte("\033[97m")
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = nil
		Red = nil
		Green = nil
		Yellow = nil
		Blue = nil
		Purple = nil
		Cyan = nil
		Gray = nil
		White = nil
	}
}
