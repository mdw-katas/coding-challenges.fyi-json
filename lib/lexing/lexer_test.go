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
	this.assertLexed(`null`, token(lexing.TokenNull, "null"))
	this.assertLexed(`true`, token(lexing.TokenTrue, "true"))
	this.assertLexed(`false`, token(lexing.TokenFalse, "false"))
	this.assertLexed(`null--trailing-bad-stuff-will-be-identified-by-parser`,
		token(lexing.TokenNull, "null"),
	)
}
func (this *Suite) TestTopLevel_Number() {
	this.assertLexed(`0`, token(lexing.TokenZero, "0"))
	this.assertLexed(`1`, token(lexing.TokenOne, "1"))
	this.assertLexed(`1234567890`,
		token(lexing.TokenOne, "1"),
		token(lexing.TokenTwo, "2"),
		token(lexing.TokenThree, "3"),
		token(lexing.TokenFour, "4"),
		token(lexing.TokenFive, "5"),
		token(lexing.TokenSix, "6"),
		token(lexing.TokenSeven, "7"),
		token(lexing.TokenEight, "8"),
		token(lexing.TokenNine, "9"),
		token(lexing.TokenZero, "0"),
	)
	this.assertLexed(`0.NaN`, token(lexing.TokenZero, "0"))
	this.assertLexed(`0.0`,
		token(lexing.TokenZero, "0"),
		token(lexing.TokenDecimalPoint, "."),
		token(lexing.TokenZero, "0"),
	)
	this.assertLexed(`0.0123456789`,
		token(lexing.TokenZero, "0"),
		token(lexing.TokenDecimalPoint, "."),
		token(lexing.TokenZero, "0"),
		token(lexing.TokenOne, "1"),
		token(lexing.TokenTwo, "2"),
		token(lexing.TokenThree, "3"),
		token(lexing.TokenFour, "4"),
		token(lexing.TokenFive, "5"),
		token(lexing.TokenSix, "6"),
		token(lexing.TokenSeven, "7"),
		token(lexing.TokenEight, "8"),
		token(lexing.TokenNine, "9"),
	)
	this.assertLexed(`1234567890.0123456789`,
		token(lexing.TokenOne, "1"),
		token(lexing.TokenTwo, "2"),
		token(lexing.TokenThree, "3"),
		token(lexing.TokenFour, "4"),
		token(lexing.TokenFive, "5"),
		token(lexing.TokenSix, "6"),
		token(lexing.TokenSeven, "7"),
		token(lexing.TokenEight, "8"),
		token(lexing.TokenNine, "9"),
		token(lexing.TokenZero, "0"),
		token(lexing.TokenDecimalPoint, "."),
		token(lexing.TokenZero, "0"),
		token(lexing.TokenOne, "1"),
		token(lexing.TokenTwo, "2"),
		token(lexing.TokenThree, "3"),
		token(lexing.TokenFour, "4"),
		token(lexing.TokenFive, "5"),
		token(lexing.TokenSix, "6"),
		token(lexing.TokenSeven, "7"),
		token(lexing.TokenEight, "8"),
		token(lexing.TokenNine, "9"),
	)
	this.assertLexed(`-1`, token(lexing.TokenNegativeSign, "-"), token(lexing.TokenOne, "1"))
	this.assertLexed(`-0`, token(lexing.TokenNegativeSign, "-"), token(lexing.TokenZero, "0"))
	this.assertLexed(`-0.1`,
		token(lexing.TokenNegativeSign, "-"),
		token(lexing.TokenZero, "0"),
		token(lexing.TokenDecimalPoint, "."),
		token(lexing.TokenOne, "1"),
	)
}

func token(tokenType lexing.TokenType, value string) lexing.Token {
	return lexing.Token{Type: tokenType, Value: []byte(value)}
}

//this.assertLexed(`""`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`[]`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`{}`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
