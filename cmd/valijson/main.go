package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/lexing"
	"github.com/mdwhatcott/coding-challenges.fyi-json/lib/printing"
)

var Version = "dev"

func main() {
	var format string
	log.SetFlags(log.Lshortfile | log.Ltime)
	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", filepath.Base(os.Args[0]), Version), flag.ExitOnError)
	flags.StringVar(&format, "fmt", "pretty", "How to format the output, one of 'pretty', 'compact', 'verbatim'.")
	flags.Usage = func() {
		_, _ = fmt.Fprintf(flags.Output(), "Usage of %s:\n", flags.Name())
		_, _ = fmt.Fprintf(flags.Output(), "%s [args ...]\n", filepath.Base(os.Args[0]))
		_, _ = fmt.Fprintln(flags.Output(), "More details here.")
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])
	if format != "pretty" && format != "compact" && format != "verbatim" {
		log.Fatalln("Invalid output format:", format)
	}

	printer := printing.NewPrettyPrinter(os.Stdout)
	if format == "compact" {
		printer = printing.NewCompactPrinter(os.Stdout)
	} else if format == "verbatim" {
		printer = printing.NewVerbatimPrinter(os.Stdout)
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalln(err)
	}
	byteCount := 0
	tokenCount := 0
	for token := range lexing.Lex(input) {
		printer.Print(token)
		if token.Type == lexing.TokenIllegal {
			log.Fatalf("Illegal token at index %d: %s", byteCount, token.Value)
		}
		tokenCount++
		byteCount += len(token.Value)
	}
	fmt.Println()
	log.Printf("JSON document with %d bytes and %d tokens validated successfully.", byteCount, tokenCount)
}
