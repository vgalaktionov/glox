package main

import (
	"fmt"
	"os"

	"github.com/vgalaktionov/glox/lox"
)

func main() {
	fmt.Println("Welcome to glox!")

	if len(os.Args) > 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println(`Usage:
	glox [script]	# executes source file
	glox		# launches repl`)
	} else if len(os.Args) == 2 {
		lox.RunFile(os.Args[1])
	} else {
		lox.RunPrompt()
	}
}
