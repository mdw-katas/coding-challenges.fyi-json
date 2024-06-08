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
	go lexer.Lex()
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
func (this *Suite) SkipTestTopLevel_Number_0() {
	this.assertSuccess(`0`, lexing.Token{Type: lexing.TokenZero, Value: []byte("0")})
}

//this.assertSuccess(`""`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertSuccess(`[]`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
//this.assertSuccess(`{}`, lexing.Token{Type: lexing.TokenFalse, Value: []byte("false")})
