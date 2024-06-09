package lexing

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"
)

var (
	ErrUnexpectedEOF        = errors.New("unexpected EOF")
	ErrUnexpectedWhitespace = errors.New("unexpected whitespace")
)

type TokenType string

const (
	TokenNull         TokenType = "<null>"
	TokenTrue         TokenType = "<true>"
	TokenFalse        TokenType = "<false>"
	TokenZero         TokenType = "<0>"
	TokenDecimalPoint TokenType = "<.>"
)

type Token struct {
	Type  TokenType
	Value []byte
}

type stateFn func(*Lexer) stateFn

type Lexer struct {
	input  []byte
	start  int // start position of this item.
	pos    int // current position in the input.
	width  int // width of last rune read from input.
	err    error
	output chan Token
}

func New(input []byte) *Lexer { // TODO: accept io.Reader?
	return &Lexer{
		input:  input,
		output: make(chan Token),
	}
}

func (this *Lexer) Output() <-chan Token {
	return this.output
}

func (this *Lexer) Error() error {
	return this.err
}

func (this *Lexer) Lex() {
	defer close(this.output)
	if len(this.input) == 0 {
		this.err = ErrUnexpectedEOF
		return
	}

	for state := lexTopLevelValue; state != nil && this.pos < len(this.input); {
		state = state(this)
	}
}

func (this *Lexer) emit(tokenType TokenType) {
	this.output <- Token{
		Type:  tokenType,
		Value: this.input[this.start:this.pos],
	}
	this.start = this.pos
}

func (this *Lexer) at(offset int) rune {
	return rune(this.input[this.pos+offset])
}

func lexTopLevelValue(this *Lexer) stateFn {
	if len(this.input) == 0 {
		this.err = fmt.Errorf("%w after %d bytes", ErrUnexpectedEOF, len(this.input))
		return nil
	}
	if unicode.IsSpace(this.at(0)) { // TODO: only consider certain low/ascii space values
		this.err = fmt.Errorf("%w at index %d", ErrUnexpectedWhitespace, 0)
		return nil
	}
	if bytes.HasPrefix(this.input, _null) {
		this.pos += len(_null)
		this.emit(TokenNull)
		return nil
	}
	if bytes.HasPrefix(this.input, _true) {
		this.pos += len(_true)
		this.emit(TokenTrue)
		return nil
	}
	if bytes.HasPrefix(this.input, _false) {
		this.pos += len(_false)
		this.emit(TokenFalse)
		return nil
	}
	if this.input[this.start] == _0 {
		return lexZero
	}
	return nil
}

func lexZero(this *Lexer) stateFn {
	this.pos++
	this.emit(TokenZero)
	return lexFraction
}

func lexFraction(this *Lexer) stateFn {
	if this.at(0) == '.' && unicode.IsDigit(this.at(1)) { // TODO: hmm, only consider ascii digits
		this.pos++
		this.emit(TokenDecimalPoint)
		this.pos++
		this.emit(TokenZero)
	}
	return nil
}

var (
	_null  = []byte("null")
	_true  = []byte("true")
	_false = []byte("false")
)

const (
	_0 = '0'
)
