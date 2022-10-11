package lox

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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

	for _, token := range tokens {
		fmt.Println(token)
	}
}

type DefaultErrorReporter struct{}

func (er *DefaultErrorReporter) Error(line int, message string, params ...interface{}) {
	er.report(line, "", fmt.Sprintf(message, params...))
}

func (er *DefaultErrorReporter) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "❌ [line %d] Error %s: %s\n", line, where, message)
	hadError = true
}
