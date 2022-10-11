package main

import (
	"fmt"
	"os"

	"github.com/vgalaktionov/glox/lox"
)

func main() {
	fmt.Println("Welcome to glox!")

	if len(os.Args) > 1 || os.Args[0] == "-h" || os.Args[0] == "--help" {
		fmt.Println(`Usage:
	glox [script]	# executes source file
	glox		# launches repl`)
	} else if len(os.Args) == 1 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
