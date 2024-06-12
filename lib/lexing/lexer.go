package lexing

import "slices"

type TokenType string

const (
	TokenIllegal     TokenType = "<ILLEGAL>"
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

type Lexer struct {
	input  []byte
	start  int
	stop   int
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

	this.lexValue()

	if this.stop < len(this.input) {
		this.emit(TokenIllegal)
	}
}

func (this *Lexer) peek() rune {
	return this.at(0)
}
func (this *Lexer) at(offset int) rune {
	if this.stop >= len(this.input) {
		return 0
	}
	return rune(this.input[this.stop+offset])
}
func (this *Lexer) stepN(n int) {
	this.stop += n
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
func (this *Lexer) acceptN(n int, set ...rune) bool {
	for x := 0; x < n; x++ {
		if !this.accept(set...) {
			return false
		}
	}
	return true
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
	if tokenType == TokenIllegal {
		this.stop = len(this.input)
	}
	this.output <- Token{Type: tokenType, Value: this.input[this.start:this.stop]}
	this.start = this.stop
}

func (this *Lexer) lexValue() bool {
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
	return true
}

func (this *Lexer) acceptNumber() bool {
	if !couldBeNumber(this.peek()) {
		return false
	}
	this.accept(sign...)
	if !isDigit(this.peek()) {
		this.stop = this.start
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
func (this *Lexer) acceptString() bool {
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
	this.emit(TokenIllegal)
	return false
}
func (this *Lexer) acceptArray() bool {
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
	this.emit(TokenIllegal)
	return false
}
func (this *Lexer) acceptObject() bool {
	if !this.accept(leftCurly) {
		return false
	}
	this.emit(TokenObjectStart)
	if this.accept(rightCurly) {
		return true
	}

	for {
		if !this.acceptString() {
			this.emit(TokenIllegal)
			return false
		} else {
			this.emit(TokenString)
		}

		if !this.accept(colon) {
			this.emit(TokenIllegal)
			return false
		} else {
			this.emit(TokenColon)
		}

		if !this.lexValue() {
			this.emit(TokenIllegal)
			return false
		}

		if !this.accept(comma) {
			break
		}
		this.emit(TokenComma)
	}

	if this.accept(rightCurly) {
		return true
	}
	this.emit(TokenIllegal)
	return false
}

func isWhiteSpace(r rune) bool  { return r == space } // TODO: additional whitespace characters
func couldBeNumber(r rune) bool { return isSign(r) || isDigit(r) }
func isDigit(r rune) bool       { return zero <= r && r <= nine }
func isSign(r rune) bool        { return r == positive || r == negative }

var (
	_null     = []rune("null")
	_true     = []rune("true")
	_false    = []rune("false")
	digits    = []rune("0123456789")
	hexDigits = append(digits, []rune("abcdefg"+"ABCDEFG")...)
	sign      = []rune{positive, negative}
	exponent  = []rune{_exponent, _Exponent}
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
