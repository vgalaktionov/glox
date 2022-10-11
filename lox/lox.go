package lox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/vgalaktionov/glox/parser"
	"github.com/vgalaktionov/glox/scanner"
)

var hadError = false

func RunFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("Failed to read source file at path: ", path)
	}
	run(string(bytes))
	if hadError {
		os.Exit(65)
	}

}

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		hasInput := scanner.Scan()
		line := fmt.Sprintln(scanner.Text())

		if len(line) == 0 || !hasInput {
			break
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
		}
		run(line)
		hadError = false
	}
}

func run(source string) {
	s := scanner.NewScanner(source, &DefaultErrorReporter{})
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	expr := p.Parse()

	// stop on syntax error
	if hadError {
		return
	}

	fmt.Println(parser.AstPrinter(expr))
}

type DefaultErrorReporter struct{}

func (er *DefaultErrorReporter) Error(context interface{}, message string, params ...interface{}) {
	switch c := context.(type) {
	case scanner.Token:
		{
			if c.Type == scanner.EOF {
				er.report(c.Line, " at end", fmt.Sprintf(message, params...))
			} else {
				er.report(c.Line, fmt.Sprintf(" at '%s'", c.Lexeme), fmt.Sprintf(message, params...))
			}
		}
	case int: // line number
		er.report(c, "", fmt.Sprintf(message, params...))
	default:
		er.report(-1, "", fmt.Sprintf(message, params...))
	}
}

func (er *DefaultErrorReporter) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "❌ [line %d] Error %s: %s\n", line, where, message)
	hadError = true
}
