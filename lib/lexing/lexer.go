package lexing

import (
	"bytes"
	"errors"
	"fmt"
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

type stateMethod func() stateMethod

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
		this.err = fmt.Errorf("%w after %d bytes", ErrUnexpectedEOF, len(this.input))
		return
	}
	if isWhiteSpace(this.at(0)) {
		this.err = fmt.Errorf("%w at index %d", ErrUnexpectedWhitespace, 0)
		return
	}

	for state := this.lexValue; state != nil && this.pos < len(this.input); {
		state = state()
	}
}

func (this *Lexer) at(offset int) rune {
	return rune(this.input[this.pos+offset])
}
func (this *Lexer) emit(tokenType TokenType) {
	this.output <- Token{
		Type:  tokenType,
		Value: this.input[this.start:this.pos],
	}
	this.start = this.pos
}

func (this *Lexer) lexValue() stateMethod {
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
		return this.lexZero
	}
	return nil
}
func (this *Lexer) lexZero() stateMethod {
	this.pos++
	this.emit(TokenZero)
	return this.lexFraction
}
func (this *Lexer) lexFraction() stateMethod {
	if this.at(0) == '.' && isDigit(this.at(1)) {
		this.pos++
		this.emit(TokenDecimalPoint)
		this.pos++
		this.emit(TokenZero)
	}
	return nil
}

func isWhiteSpace(r rune) bool {
	return r == ' ' // TODO: additional whitespace characters
}
func isDigit(r rune) bool {
	switch r {
	case '0' /*, '1', '2', '3', '4', '5', '6', '7', '8', '9' */ :
		return true
	}
	return false
}

var (
	_null  = []byte("null")
	_true  = []byte("true")
	_false = []byte("false")
)

const (
	_0 = '0'
)
