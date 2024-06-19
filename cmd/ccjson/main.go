package main

import (
	"bytes"
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

const exampleInput = `{"foo":"bar","baz":[1,2,3]}`

func main() {
	var format string
	log.SetFlags(0)
	log.SetPrefix("[LOG] ")
	program := filepath.Base(os.Args[0])
	flags := flag.NewFlagSet(fmt.Sprintf("%s @ %s", program, Version), flag.ExitOnError)
	flags.StringVar(&format, "fmt", "colors", "How to format the output, one of 'colors', 'pretty', 'compact', 'verbatim'.")
	flags.Usage = func() {
		_, _ = fmt.Fprintln(flags.Output(), flags.Name())
		_, _ = fmt.Fprintln(flags.Output(), "> Validates JSON data from stdin, outputs JSON to stdout.")
		_, _ = fmt.Fprintln(flags.Output(), "> Example usage:")
		_, _ = fmt.Fprintf(flags.Output(), `$ echo -n '%s' | %s -fmt pretty`+"\n", exampleInput, program)
		validateJSON(flags.Output(), bytes.NewBufferString(exampleInput), "pretty")
		_, _ = fmt.Fprintln(flags.Output(), "> Flags:")
		flags.PrintDefaults()
	}
	_ = flags.Parse(os.Args[1:])
	if format != "colors" && format != "pretty" && format != "compact" && format != "verbatim" {
		log.Fatalln("Invalid output format:", format)
	}

	validateJSON(os.Stdout, os.Stdin, format)
}
func validateJSON(output io.Writer, input io.Reader, format string) {
	byteCount := 0
	tokenCount := 0
	printer := newPrinter(output, format)
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
func newPrinter(output io.Writer, format string) printing.Printer {
	switch format {
	case "colors":
		return printing.NewColorPrinter(output, printing.NewPrettyPrinter(output))
	case "pretty":
		return printing.NewPrettyPrinter(output)
	case "compact":
		return printing.NewCompactPrinter(output)
	case "verbatim":
		return printing.NewVerbatimPrinter(output)
	default:
		panic("invalid format: " + format)
	}
}
