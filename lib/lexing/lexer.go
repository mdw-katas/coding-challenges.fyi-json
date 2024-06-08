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
	TokenWhitespace TokenType = "<whitespace>"
	TokenNull       TokenType = "<null>"
	TokenTrue       TokenType = "<true>"
	TokenFalse      TokenType = "<false>"
	TokenZero       TokenType = "<0>"
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

func New(input []byte) *Lexer {
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
	for state := lexTopLevelValue; state != nil; {
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

func lexTopLevelValue(this *Lexer) stateFn {
	if len(this.input) == 0 {
		this.err = fmt.Errorf("%w after %d bytes", ErrUnexpectedEOF, len(this.input))
		return nil
	}
	if unicode.IsSpace(rune(this.input[0])) {
		this.err = fmt.Errorf("%w at index %d", ErrUnexpectedWhitespace, 0)
		return nil
	}
	if bytes.HasPrefix(this.input, []byte("null")) {
		this.pos += len("null")
		this.emit(TokenNull)
		return nil
	}
	if bytes.HasPrefix(this.input, []byte("true")) {
		this.pos += len("true")
		this.emit(TokenTrue)
		return nil
	}
	if bytes.HasPrefix(this.input, []byte("false")) {
		this.pos += len("false")
		this.emit(TokenFalse)
		return nil
	}
	if this.input[this.start] == '0' {
		return lexZero(this)
	}
	return nil
}

func lexZero(this *Lexer) stateFn {
	this.pos++
	this.emit(TokenZero)
	return lexFraction(this)
}

func lexFraction(this *Lexer) stateFn {
	return nil
}
