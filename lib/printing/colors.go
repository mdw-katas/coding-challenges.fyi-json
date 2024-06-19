package printing

import (
	"io"

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
	switch token.Type {
	case lexing.TokenNull:
		this.write(gray, token)
	case lexing.TokenTrue:
		this.write(green, token)
	case lexing.TokenFalse:
		this.write(purple, token)
	case lexing.TokenNumber:
		this.write(yellow, token)
	case lexing.TokenString:
		this.write(blue, token)
	case lexing.TokenArrayStart,
		lexing.TokenArrayStop,
		lexing.TokenObjectStart,
		lexing.TokenObjectStop,
		lexing.TokenComma,
		lexing.TokenColon:
		this.write(cyan, token)
	case lexing.TokenIllegal:
		this.write(red, token)
	default:
		this.inner.Print(token)
	}
}

func (this *colors) write(color []byte, token lexing.Token) {
	_, _ = this.out.Write(color)
	this.inner.Print(token)
	_, _ = this.out.Write(reset)
}

var (
	reset  = []byte("\033[0m")
	red    = []byte("\033[31m")
	green  = []byte("\033[32m")
	yellow = []byte("\033[33m")
	blue   = []byte("\033[34m")
	purple = []byte("\033[35m")
	cyan   = []byte("\033[36m")
	gray   = []byte("\033[37m")
	white  = []byte("\033[97m")
)
