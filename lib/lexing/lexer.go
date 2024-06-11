package lexing

import "slices"

type TokenType string

const (
	TokenIllegal TokenType = "<ILLEGAL>"
	TokenNull    TokenType = "<null>"
	TokenTrue    TokenType = "<true>"
	TokenFalse   TokenType = "<false>"
	TokenNumber  TokenType = "<number>"
)

type Token struct {
	Type  TokenType
	Value []byte
}

type stateMethod func() stateMethod

type Lexer struct {
	input  []byte
	start  int
	pos    int
	output chan Token
}

func Lex(input []byte) chan Token { // TODO: accept io.Reader?
	lexer := &Lexer{input: input, output: make(chan Token)}
	go lexer.lex()
	return lexer.output
}
func (this *Lexer) lex() {
	defer close(this.output)
	if len(this.input) == 0 {
		return
	}
	if isWhiteSpace(this.peek()) {
		return
	}
	for state := this.lexValue; state != nil && this.pos < len(this.input); {
		state = state()
	}
}

func (this *Lexer) peek() rune {
	return this.at(0)
}
func (this *Lexer) at(offset int) rune {
	if this.pos >= len(this.input) {
		return 0
	}
	return rune(this.input[this.pos+offset])
}
func (this *Lexer) from(offset int) rune {
	return rune(this.input[this.start+offset])
}
func (this *Lexer) stepN(n int) {
	this.pos += n
}
func (this *Lexer) step() {
	this.stepN(1)
}
func (this *Lexer) accept(set ...rune) bool {
	ok := slices.Index(set, this.peek()) >= 0
	if ok {
		this.step()
	}
	return ok
}
func (this *Lexer) acceptRun(set ...rune) (result int) {
	for {
		if !this.accept(set...) {
			return result
		}
		result++
	}
}
func (this *Lexer) acceptSequence(sequence []rune) bool {
	for _, s := range sequence {
		if !this.accept(s) {
			return false
		}
	}
	return true
}
func (this *Lexer) emit(tokenType TokenType) {
	this.output <- Token{Type: tokenType, Value: this.input[this.start:this.pos]}
	this.start = this.pos
}

func (this *Lexer) lexValue() stateMethod {
	if this.acceptSequence(_null) {
		this.emit(TokenNull)
	} else if this.acceptSequence(_true) {
		this.emit(TokenTrue)
	} else if this.acceptSequence(_false) {
		this.emit(TokenFalse)
	} else if couldBeNumber(this.peek()) {
		if this.acceptNumber() {
			this.emit(TokenNumber)
		} else {
			this.emit(TokenIllegal)
		}
	}
	return nil
}
func (this *Lexer) acceptNumber() bool {
	this.accept(sign...)
	if !isDigit(this.peek()) {
		this.pos = this.start
		return false
	}
	if !this.accept(zero) {
		if this.acceptRun(digits...) == 0 {
			return false
		}
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

func isWhiteSpace(r rune) bool  { return r == ' ' } // TODO: additional whitespace characters
func couldBeNumber(r rune) bool { return isSign(r) || isDigit(r) }
func isDigit(r rune) bool       { return zero <= r && r <= nine }
func isSign(r rune) bool        { return r == positive || r == negative }

var (
	_null    = []rune("null")
	_true    = []rune("true")
	_false   = []rune("false")
	digits   = []rune("0123456789")
	sign     = []rune{positive, negative}
	exponent = []rune{_exponent, _Exponent}
)

const (
	positive     = '+'
	negative     = '-'
	_exponent    = 'e'
	_Exponent    = 'E'
	decimalPoint = '.'
	zero         = '0'
	nine         = '9'
)
