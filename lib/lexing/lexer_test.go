package lexing

import (
	"fmt"
	"testing"

	"github.com/mdwhatcott/testing/should"
)

func TestLex(t *testing.T) {
	t.Run("top-level", func(t *testing.T) {
		testLex(t, "")
		testLex(t, " ")
		testLex(t, `null`, token(TokenNull, "null"))
		testLex(t, `true`, token(TokenTrue, "true"))
		testLex(t, `false`, token(TokenFalse, "false"))
		testLex(t, `null--trailing-bad-stuff`,
			token(TokenNull, "null"),
			token(TokenIllegal, "--trailing-bad-stuff"),
		)
	})
	t.Run("numbers", func(t *testing.T) {
		testLex(t, `0`, token(TokenNumber, "0"))
		testLex(t, `0a`, token(TokenNumber, "0"), token(TokenIllegal, "a"))
		testLex(t, `1`, token(TokenNumber, "1"))
		testLex(t, `1234567890`, token(TokenNumber, "1234567890"))
		testLex(t, `0.NaN`, token(TokenIllegal, "0.NaN"))
		testLex(t, `0.0`, token(TokenNumber, "0.0"))
		testLex(t, `0.0123456789`, token(TokenNumber, "0.0123456789"))
		testLex(t, `1234567890.0123456789`, token(TokenNumber, "1234567890.0123456789"))
		testLex(t, `-1`, token(TokenNumber, "-1"))
		testLex(t, `-0`, token(TokenNumber, "-0"))
		testLex(t, `-0.1`, token(TokenNumber, "-0.1"))
		testLex(t, `3.7e-5`, token(TokenNumber, "3.7e-5"))
		testLex(t, `3.7e+5`, token(TokenNumber, "3.7e+5"))
	})
	t.Run("strings", func(t *testing.T) {
		testLex(t, `"`, token(TokenIllegal, `"`))
		testLex(t, `""`, token(TokenString, `""`))
		testLex(t, `"a"`, token(TokenString, `"a"`))
		testLex(t, `"ab"`, token(TokenString, `"ab"`))
		testLex(t, `"\t"`, token(TokenString, `"\t"`))
		testLex(t, `"\r"`, token(TokenString, `"\r"`))
		testLex(t, `"\n"`, token(TokenString, `"\n"`))
		testLex(t, `"\f"`, token(TokenString, `"\f"`))
		testLex(t, `"\b"`, token(TokenString, `"\b"`))
		testLex(t, `"\\"`, token(TokenString, `"\\"`))
		testLex(t, `"\/"`, token(TokenString, `"\/"`))
		testLex(t, `"\""`, token(TokenString, `"\""`))
		testLex(t, `"\u1234"`, token(TokenString, `"\u1234"`))
		testLex(t, `"\u12aB"`, token(TokenString, `"\u12aB"`))
		testLex(t, `"\uaB12"`, token(TokenString, `"\uaB12"`))
		testLex(t, `"\u5678"`, token(TokenString, `"\u5678"`))
		testLex(t, `"\u90ab"`, token(TokenString, `"\u90ab"`))
		testLex(t, `"\ucdef"`, token(TokenString, `"\ucdef"`))
		testLex(t, `"\uCDEF"`, token(TokenString, `"\uCDEF"`))
		testLex(t, `"\u123x"`, token(TokenIllegal, `"\u123x"`))
		testLex(t, `"`+"\t"+`"`, token(TokenIllegal, `"`+"\t"+`"`))
	})
	t.Run("arrays", func(t *testing.T) {
		testLex(t, `[`, token(TokenArrayStart, `[`), token(TokenIllegal, ""))
		testLex(t, `[1`, token(TokenArrayStart, `[`), token(TokenNumber, `1`), token(TokenIllegal, ""))
		testLex(t, `[]`, token(TokenArrayStart, `[`), token(TokenArrayStop, `]`))
		testLex(t, `[1]`,
			token(TokenArrayStart, `[`),
			token(TokenNumber, "1"),
			token(TokenArrayStop, `]`),
		)
		testLex(t, `[1,2]`,
			token(TokenArrayStart, `[`),
			token(TokenNumber, "1"),
			token(TokenComma, ","),
			token(TokenNumber, "2"),
			token(TokenArrayStop, "]"),
		)
		testLex(t, `[1,"b",3]`,
			token(TokenArrayStart, `[`),
			token(TokenNumber, "1"),
			token(TokenComma, ","),
			token(TokenString, `"b"`),
			token(TokenComma, ","),
			token(TokenNumber, "3"),
			token(TokenArrayStop, "]"),
		)
		testLex(t, `[1,["b"],3]`,
			token(TokenArrayStart, `[`),
			token(TokenNumber, `1`),
			token(TokenComma, `,`),
			token(TokenArrayStart, `[`),
			token(TokenString, `"b"`),
			token(TokenArrayStop, `]`),
			token(TokenComma, `,`),
			token(TokenNumber, `3`),
			token(TokenArrayStop, `]`),
		)
	})
}
func lex(s string) (result []Token) {
	defer func() { recover() }()
	for token := range Lex([]byte(s)) {
		result = append(result, token)
	}
	return result
}
func testLex(t *testing.T, input string, expected ...Token) {
	t.Run(input, func(t *testing.T) {
		should.So(t, lex(input), should.Equal, expected)
	})
}
func token(tokenType TokenType, value string) Token {
	return Token{Type: tokenType, Value: []byte(value)}
}
func (this Token) GoString() string {
	return fmt.Sprintf(`lexing.Token{Type:"%s", Value: []byte("%s")}`, this.Type, this.Value)
}
