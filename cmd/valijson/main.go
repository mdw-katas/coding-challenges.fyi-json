package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
)

var Version = "dev"

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", filepath.Base(os.Args[0]), Version), flag.ExitOnError)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage of %s:\n", flags.Name())
		_, _ = fmt.Fprintf(flags.Output(), "%s [args ...]\n", filepath.Base(os.Args[0]))
		_, _ = fmt.Fprintln(flags.Output(), "More details here.")
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	byteCount := 0
	tokenCount := 0
	for token := range lexing.Lex(input) {
		if token.Type == lexing.TokenIllegal {
			log.Fatalf("Illegal token at index %d: %s", byteCount, token.Value)
		}
		// TODO: feed token to a printer
		tokenCount++
		byteCount += len(token.Value)
	}
	log.Printf("JSON document with %d bytes and %d tokens validated successfully.", byteCount, tokenCount)
}
