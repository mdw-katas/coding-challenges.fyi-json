package lexing

import (
	"errors"
	"fmt"
)

var (
	ErrUnexpectedEOF        = errors.New("unexpected EOF")
	ErrUnexpectedWhitespace = errors.New("unexpected whitespace")
)

type TokenType string

const (
	TokenWhitespace TokenType = "<whitespace>"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
}

func New() *Lexer {
	return &Lexer{}
}

func (this *Lexer) Lex(raw []byte) (result []Token, err error) {
	if len(raw) == 0 {
		err = fmt.Errorf("%w after %d bytes", ErrUnexpectedEOF, len(raw))
	} else {
		err = fmt.Errorf("%w at index %d", ErrUnexpectedWhitespace, 0)
	}
	return result, err
}
