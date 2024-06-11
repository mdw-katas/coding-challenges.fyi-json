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

func (this *Suite) lex(s string) []lexing.Token {
	lexer := lexing.New([]byte(s))
	go func() {
		defer func() { recover() }()
		lexer.Lex()
	}()
	var result []lexing.Token
	for token := range lexer.Output() {
		result = append(result, token)
	}
	return result
}
func (this *Suite) assertLexed(input string, expected ...lexing.Token) {
	this.So(this.lex(input), should.Equal, expected)
}

func (this *Suite) TestTopLevelValues() {
	this.assertLexed("")
	this.assertLexed(" ")
	this.assertLexed(`null`, lexing.Token{Type: lexing.TokenNull, Value: []byte("null")})
	this.assertLexed(`true`, lexing.Token{Type: lexing.TokenTrue, Value: []byte("true")})
	this.assertLexed(`false`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
	this.assertLexed(`null--trailing-bad-stuff-will-be-identified-by-parser`,
		lexing.Token{Type: lexing.TokenNull, Value: []byte("null")},
	)
}
func (this *Suite) TestTopLevel_Number() {
	this.assertLexed(`0`, lexing.Token{Type: lexing.TokenZero, Value: []byte("0")})
	this.assertLexed(`1`, lexing.Token{Type: lexing.TokenOne, Value: []byte("1")})
	this.assertLexed(`1234567890`,
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
	this.assertLexed(`0.NaN`, lexing.Token{Type: lexing.TokenZero, Value: []byte("0")})
	this.assertLexed(`0.0`,
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
		lexing.Token{Type: lexing.TokenDecimalPoint, Value: []byte(".")},
		lexing.Token{Type: lexing.TokenZero, Value: []byte("0")},
	)
	this.assertLexed(`0.0123456789`,
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
	this.assertLexed(`1234567890.0123456789`,
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

//this.assertLexed(`""`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`[]`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`{}`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
