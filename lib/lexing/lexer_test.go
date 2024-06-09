package lexing_test

import (
	"testing"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/testing/should"
)

func TestSuite(t *testing.T) {
	should.Run(&Suite{T: should.New(t)}, should.Options.UnitTests())
}

type Suite struct {
	*should.T
}

func (this *Suite) lex(s string) ([]lexing.Token, error) {
	lexer := lexing.New([]byte(s))
	go func() {
		defer func() { recover() }()
		lexer.Lex()
	}()
	output := lexer.Output()
	var result []lexing.Token
	for token := range output {
		result = append(result, token)
	}
	return result, lexer.Error()
}
func (this *Suite) assertSuccess(input string, expected ...lexing.Token) {
	tokens, err := this.lex(input)
	this.So(err, should.BeNil)
	this.So(tokens, should.Equal, expected)
}
func (this *Suite) assertFailure(input string, expected error) {
	tokens, err := this.lex(input)
	this.So(tokens, should.BeNil)
	this.So(err, should.WrapError, expected)
}

func (this *Suite) TestTopLevel_Blank() {
	this.assertFailure("", lexing.ErrUnexpectedEOF)
}
func (this *Suite) TestTopLevel_JustWhitespace() {
	this.assertFailure(" ", lexing.ErrUnexpectedWhitespace)
}
func (this *Suite) TestTopLevel_Null() {
	this.assertSuccess(`null`, lexing.Token{Type: lexing.TokenNull, Value: []byte("null")})
}
func (this *Suite) TestTopLevel_True() {
	this.assertSuccess(`true`, lexing.Token{Type: lexing.TokenTrue, Value: []byte("true")})
}
func (this *Suite) TestTopLevel_False() {
	this.assertSuccess(`false`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
}
func (this *Suite) TestTopLevel_Number() {
	this.assertSuccess(`0`, lexing.Token{Type: lexing.TokenZero, Value: []byte("0")})
	this.assertSuccess(`1`, lexing.Token{Type: lexing.TokenOne, Value: []byte("1")})
	this.assertSuccess(`1234567890`,
		lexing.Token{Type: lexing.TokenOne, Value: []byte("1")},
		lexing.Token{Type: lexing.TokenTwo, Value: []byte("2")},
		lexing.Token{Type: lexing.TokenThree, Value: []byte("3")},
		lexing.Token{Type: lexing.TokenFour, Value: []byte("4")},
		lexing.Token{Type: lexing.TokenFive, Value: []byte("5")},
		lexing.Token{Type: lexing.TokenSix, Value: []byte("6")},
		lexing.Token{Type: lexing.TokenSeven, Value: []byte("7")},
		lexing.Token{Type: lexing.TokenEight, Value: []byte("8")},
		lexing.Token{Type: lexing.TokenNine, Value: []byte("9")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
	)
	this.assertSuccess(`0.0`,
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenDecimalPoint, Value: []byte(".")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
	)
	this.assertSuccess(`0.0123456789`,
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenDecimalPoint, Value: []byte(".")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenOne, Value: []byte("1")},
		lexing.Token{Type: lexing.TokenTwo, Value: []byte("2")},
		lexing.Token{Type: lexing.TokenThree, Value: []byte("3")},
		lexing.Token{Type: lexing.TokenFour, Value: []byte("4")},
		lexing.Token{Type: lexing.TokenFive, Value: []byte("5")},
		lexing.Token{Type: lexing.TokenSix, Value: []byte("6")},
		lexing.Token{Type: lexing.TokenSeven, Value: []byte("7")},
		lexing.Token{Type: lexing.TokenEight, Value: []byte("8")},
		lexing.Token{Type: lexing.TokenNine, Value: []byte("9")},
	)
	this.assertSuccess(`1234567890.0123456789`,
		lexing.Token{Type: lexing.TokenOne, Value: []byte("1")},
		lexing.Token{Type: lexing.TokenTwo, Value: []byte("2")},
		lexing.Token{Type: lexing.TokenThree, Value: []byte("3")},
		lexing.Token{Type: lexing.TokenFour, Value: []byte("4")},
		lexing.Token{Type: lexing.TokenFive, Value: []byte("5")},
		lexing.Token{Type: lexing.TokenSix, Value: []byte("6")},
		lexing.Token{Type: lexing.TokenSeven, Value: []byte("7")},
		lexing.Token{Type: lexing.TokenEight, Value: []byte("8")},
		lexing.Token{Type: lexing.TokenNine, Value: []byte("9")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenDecimalPoint, Value: []byte(".")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenOne, Value: []byte("1")},
		lexing.Token{Type: lexing.TokenTwo, Value: []byte("2")},
		lexing.Token{Type: lexing.TokenThree, Value: []byte("3")},
		lexing.Token{Type: lexing.TokenFour, Value: []byte("4")},
		lexing.Token{Type: lexing.TokenFive, Value: []byte("5")},
		lexing.Token{Type: lexing.TokenSix, Value: []byte("6")},
		lexing.Token{Type: lexing.TokenSeven, Value: []byte("7")},
		lexing.Token{Type: lexing.TokenEight, Value: []byte("8")},
		lexing.Token{Type: lexing.TokenNine, Value: []byte("9")},
	)
}

//this.assertSuccess(`""`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertSuccess(`[]`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertSuccess(`{}`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
