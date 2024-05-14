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
	result, err := lexing.New().Lex([]byte(s))
	if result != nil {
		this.Println("tokens:", result)
	}
	if err != nil {
		this.Println("err:", err)
	}
	return result, err
}

func (this *Suite) TestNoInput() {
	tokens, err := this.lex("")
	this.So(tokens, should.BeNil)
	this.So(err, should.WrapError, lexing.ErrUnexpectedEOF)
}
func (this *Suite) TestOnlyWhitespace() {
	tokens, err := this.lex(" ")
	this.So(tokens, should.BeNil)
	this.So(err, should.WrapError, lexing.ErrUnexpectedWhitespace)
}
