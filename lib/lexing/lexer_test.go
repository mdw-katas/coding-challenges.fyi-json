package lexing

import (
	"testing"

	"github.com/mdwhatcott/testing/should"
)

func TestSuite(t *testing.T) {
	should.Run(&Suite{T: should.New(t)}, should.Options.UnitTests())
}

type Suite struct {
	*should.T
}

func (this *Suite) lex(s string) []Token {
	lexer := New([]byte(s))
	go func() {
		defer func() { recover() }()
		lexer.Lex()
	}()
	var result []Token
	for token := range lexer.Output() {
		result = append(result, token)
	}
	return result
}
func (this *Suite) assertLexed(input string, expected ...Token) {
	this.So(this.lex(input), should.Equal, expected)
}

func (this *Suite) TestTopLevelValues() {
	this.assertLexed("")
	this.assertLexed(" ")
	this.assertLexed(`null`, token(TokenNull, "null"))
	this.assertLexed(`true`, token(TokenTrue, "true"))
	this.assertLexed(`false`, token(TokenFalse, "false"))
	this.assertLexed(`null--trailing-bad-stuff-will-be-identified-by-parser`,
		token(TokenNull, "null"),
	)
}
func (this *Suite) TestTopLevel_Number() {
	this.assertLexed(`0`, token(TokenZero, "0"))
	this.assertLexed(`1`, token(TokenOne, "1"))
	this.assertLexed(`1234567890`,
		token(TokenOne, "1"),
		token(TokenTwo, "2"),
		token(TokenThree, "3"),
		token(TokenFour, "4"),
		token(TokenFive, "5"),
		token(TokenSix, "6"),
		token(TokenSeven, "7"),
		token(TokenEight, "8"),
		token(TokenNine, "9"),
		token(TokenZero, "0"),
	)
	this.assertLexed(`0.NaN`, token(TokenZero, "0"))
	this.assertLexed(`0.0`,
		token(TokenZero, "0"),
		token(TokenDecimalPoint, "."),
		token(TokenZero, "0"),
	)
	this.assertLexed(`0.0123456789`,
		token(TokenZero, "0"),
		token(TokenDecimalPoint, "."),
		token(TokenZero, "0"),
		token(TokenOne, "1"),
		token(TokenTwo, "2"),
		token(TokenThree, "3"),
		token(TokenFour, "4"),
		token(TokenFive, "5"),
		token(TokenSix, "6"),
		token(TokenSeven, "7"),
		token(TokenEight, "8"),
		token(TokenNine, "9"),
	)
	this.assertLexed(`1234567890.0123456789`,
		token(TokenOne, "1"),
		token(TokenTwo, "2"),
		token(TokenThree, "3"),
		token(TokenFour, "4"),
		token(TokenFive, "5"),
		token(TokenSix, "6"),
		token(TokenSeven, "7"),
		token(TokenEight, "8"),
		token(TokenNine, "9"),
		token(TokenZero, "0"),
		token(TokenDecimalPoint, "."),
		token(TokenZero, "0"),
		token(TokenOne, "1"),
		token(TokenTwo, "2"),
		token(TokenThree, "3"),
		token(TokenFour, "4"),
		token(TokenFive, "5"),
		token(TokenSix, "6"),
		token(TokenSeven, "7"),
		token(TokenEight, "8"),
		token(TokenNine, "9"),
	)
	this.assertLexed(`-1`,
		token(TokenNegativeSign, "-"),
		token(TokenOne, "1"),
	)
	this.assertLexed(`-0`,
		token(TokenNegativeSign, "-"),
		token(TokenZero, "0"),
	)
	this.assertLexed(`-0.1`,
		token(TokenNegativeSign, "-"),
		token(TokenZero, "0"),
		token(TokenDecimalPoint, "."),
		token(TokenOne, "1"),
	)
}

func token(tokenType TokenType, value string) Token {
	return Token{Type: tokenType, Value: []byte(value)}
}

//this.assertLexed(`""`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`[]`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertLexed(`{}`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
