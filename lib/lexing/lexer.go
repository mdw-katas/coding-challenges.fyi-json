package lexing

import (
	"io"
	"slices"
)

type TokenType string

const (
	TokenIllegal     TokenType = "<ILLEGAL>"
	TokenWhitespace  TokenType = "<whitespace>"
	TokenNull        TokenType = "<null>"
	TokenTrue        TokenType = "<true>"
	TokenFalse       TokenType = "<false>"
	TokenNumber      TokenType = "<number>"
	TokenString      TokenType = "<string>"
	TokenArrayStart  TokenType = "<[>"
	TokenArrayStop   TokenType = "<]>"
	TokenComma       TokenType = "<,>"
	TokenObjectStart TokenType = "<{>"
	TokenObjectStop  TokenType = "<}>"
	TokenColon       TokenType = "<:>"
)

type Token struct {
	Type  TokenType
	Value []byte
}

type lexer struct {
	source io.Reader
	chunk  []byte
	input  []byte // TODO: prevent input from buffering the entirety of the source data.
	start  int
	stop   int
	output chan Token
}

func Lex(source io.Reader) chan Token {
	lexer := &lexer{source: source, chunk: make([]byte, 1024), output: make(chan Token)}
	go lexer.lex()
	return lexer.output
}
func (this *lexer) lex() {
	defer close(this.output)

	this.readChunk()

	if len(this.input) == 0 {
		return
	}
	if !this.lexValue() {
		this.emit(TokenIllegal)
		return
	}
	chunk, _ := io.ReadAll(this.source)
	this.input = append(this.input, chunk...)
	if this.stop < len(this.input) {
		this.emit(TokenIllegal)
	}
}

func (this *lexer) readChunk() {
	n, _ := this.source.Read(this.chunk)
	this.input = append(this.input, this.chunk[:n]...)
	clear(this.chunk)
}

func (this *lexer) peek() rune {
	return this.at(0)
}
func (this *lexer) at(offset int) rune {
	if this.stop >= len(this.input) {
		return 0
	}
	return rune(this.input[this.stop+offset])
}
func (this *lexer) stepN(n int) {
	this.stop += n
	if this.stop >= len(this.input) {
		this.readChunk()
	}
}
func (this *lexer) step() {
	this.stepN(1)
}
func (this *lexer) accept(set ...rune) bool {
	ok := slices.Index(set, this.peek()) >= 0
	if ok {
		this.step()
	}
	return ok
}
func (this *lexer) acceptN(n int, set ...rune) bool {
	for x := 0; x < n; x++ {
		if !this.accept(set...) {
			return false
		}
	}
	return true
}
func (this *lexer) acceptRun(set ...rune) (result int) {
	for {
		if !this.accept(set...) {
			return result
		}
		result++
	}
}
func (this *lexer) acceptSequence(sequence []rune) bool {
	for _, s := range sequence {
		if !this.accept(s) {
			return false
		}
	}
	return true
}
func (this *lexer) emit(tokenType TokenType) {
	if tokenType == TokenIllegal {
		this.stop = len(this.input)
	}
	this.output <- Token{Type: tokenType, Value: this.input[this.start:this.stop]}
	this.start = this.stop
}

func (this *lexer) lexValue() bool {
	this.acceptWhitespace()

	if this.acceptSequence(_null) {
		this.emit(TokenNull)
	} else if this.acceptSequence(_true) {
		this.emit(TokenTrue)
	} else if this.acceptSequence(_false) {
		this.emit(TokenFalse)
	} else if this.acceptNumber() {
		this.emit(TokenNumber)
	} else if this.acceptString() {
		this.emit(TokenString)
	} else if this.acceptArray() {
		this.emit(TokenArrayStop)
	} else if this.acceptObject() {
		this.emit(TokenObjectStop)
	} else {
		return false
	}
	this.acceptWhitespace()
	return true
}

func (this *lexer) acceptWhitespace() {
	if this.acceptRun(whitespaces...) > 0 {
		this.emit(TokenWhitespace)
	}
}

func (this *lexer) acceptNumber() bool {
	if !couldBeNumber(this.peek()) {
		return false
	}
	this.accept(sign...)
	if !this.accept(zero) {
		this.acceptRun(digits...)
	}
	if this.accept(decimalPoint) {
		if this.acceptRun(digits...) == 0 {
			return false
		}
	}
	if this.accept(exponent...) {
		this.accept(sign...)
		if this.acceptRun(digits...) == 0 {
			return false
		}
	}
	return true
}
func (this *lexer) acceptString() bool {
	if !this.accept(quote) {
		return false
	}
	for this.stop < len(this.input) {
		switch this.peek() {
		case reverseSolidus:
			switch this.at(1) {
			case quote, reverseSolidus, solidus, backspace, formFeed, lineFeed, carriageReturn, tab:
				this.stepN(2)
			case unicode:
				this.stepN(2)
				if this.acceptN(4, hexDigits...) {
					continue
				} else {
					return false
				}
			}
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
			0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
			0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F:
			return false
		case quote:
			this.accept(quote)
			return true
		default:
			this.step()
		}
	}
	return false
}
func (this *lexer) acceptArray() bool {
	if !this.accept(leftSquare) {
		return false
	}
	this.emit(TokenArrayStart)
	if this.accept(rightSquare) {
		return true
	}
	this.lexValue()
	for {
		if this.accept(comma) {
			this.emit(TokenComma)
			this.lexValue()
		} else {
			break
		}
	}
	if this.accept(rightSquare) {
		return true
	}
	return false
}
func (this *lexer) acceptObject() bool {
	if !this.accept(leftCurly) {
		return false
	}
	this.emit(TokenObjectStart)
	this.acceptWhitespace()
	if this.accept(rightCurly) {
		return true
	}

	for {
		this.acceptWhitespace()
		if !this.acceptString() {
			return false
		} else {
			this.emit(TokenString)
		}

		this.acceptWhitespace()

		if !this.accept(colon) {
			return false
		} else {
			this.emit(TokenColon)
		}

		this.acceptWhitespace()

		if !this.lexValue() {
			return false
		}

		this.acceptWhitespace()

		if !this.accept(comma) {
			break
		}
		this.emit(TokenComma)

		this.acceptWhitespace()
	}

	this.acceptWhitespace()

	if this.accept(rightCurly) {
		return true
	}
	return false
}

func couldBeNumber(r rune) bool { return isSign(r) || isDigit(r) }
func isDigit(r rune) bool       { return zero <= r && r <= nine }
func isSign(r rune) bool        { return r == positive || r == negative }

var (
	_null       = []rune("null")
	_true       = []rune("true")
	_false      = []rune("false")
	whitespaces = []rune{space, '\n', '\r', '\t'}
	digits      = []rune("0123456789")
	hexDigits   = append(digits, []rune("abcdefg"+"ABCDEFG")...)
	sign        = []rune{positive, negative}
	exponent    = []rune{_exponent, _Exponent}
)

const (
	positive       = '+'
	negative       = '-'
	_exponent      = 'e'
	_Exponent      = 'E'
	decimalPoint   = '.'
	zero           = '0'
	nine           = '9'
	quote          = '"'
	comma          = ','
	colon          = ':'
	leftSquare     = '['
	rightSquare    = ']'
	leftCurly      = '{'
	rightCurly     = '}'
	space          = ' '
	backspace      = 'b'
	lineFeed       = 'n'
	carriageReturn = 'r'
	formFeed       = 'f'
	tab            = 't'
	unicode        = 'u'
	solidus        = '/'
	reverseSolidus = '\\'
)
