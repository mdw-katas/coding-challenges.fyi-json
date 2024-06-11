package lexing

import (
	"fmt"
	"testing"

	"github.com/mdwhatcott/testing/should"
)

func TestLex(t *testing.T) {
	testLex(t, "")
	testLex(t, " ")
	testLex(t, `null`, token(TokenNull, "null"))
	testLex(t, `true`, token(TokenTrue, "true"))
	testLex(t, `false`, token(TokenFalse, "false"))
	testLex(t, `null--trailing-bad-stuff-will-be-identified-by-parser`, token(TokenNull, "null"))
	testLex(t, `0`, token(TokenNumber, "0"))
	testLex(t, `1`, token(TokenNumber, "1"))
	testLex(t, `1234567890`, token(TokenNumber, "1234567890"))
	testLex(t, `0.NaN`, token(TokenIllegal, "0."))
	testLex(t, `0.0`, token(TokenNumber, "0.0"))
	testLex(t, `0.0123456789`, token(TokenNumber, "0.0123456789"))
	testLex(t, `1234567890.0123456789`, token(TokenNumber, "1234567890.0123456789"))
	testLex(t, `-1`, token(TokenNumber, "-1"))
	testLex(t, `-0`, token(TokenNumber, "-0"))
	testLex(t, `-0.1`, token(TokenNumber, "-0.1"))
	testLex(t, `3.7e-5`, token(TokenNumber, "3.7e-5"))
	testLex(t, `3.7e+5`, token(TokenNumber, "3.7e+5"))

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
