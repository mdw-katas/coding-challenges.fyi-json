package printing

import "github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"

type Printer interface {
	Print(lexing.Token)
}
