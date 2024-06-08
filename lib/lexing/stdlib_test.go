package lexing

import (
	"encoding/json"
	"testing"

	"github.com/mdwhatcott/testing/should"
)

func TestValidTopLevelObjects(t *testing.T) {
	assertValidJSON(t, `null`)
	assertValidJSON(t, `true`)
	assertValidJSON(t, `false`)
	assertValidJSON(t, `1`)
	assertValidJSON(t, `"a"`)
	assertValidJSON(t, `[]`)
	assertValidJSON(t, `{}`)
}
func assertValidJSON(t *testing.T, input string) {
	t.Run(input, func(t *testing.T) {
		t.Helper()
		should.So(t, json.Valid([]byte(input)), should.BeTrue)
	})
}
