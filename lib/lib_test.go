package lib

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

func (this *Suite) Setup() {}

func (this *Suite) Test() {
	this.So(1, should.Equal, 1)
}
