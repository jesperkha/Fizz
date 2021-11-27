package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	ct "github.com/daviddengcn/go-colortext"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

// Todo: Make --version option
var cmdOptions = map[string]func(){
	"help": func() {
		fmt.Println("use: fizz [filename.fizz | --option]")
	},
}

// Todo: Fix --help error msg that pops up
func RunInterpeter(args []string) {
	if len(args) == 2 {
		arg := args[1]

		// Run option commands
		if strings.HasPrefix(arg, "--") {
			option, ok := cmdOptions[strings.TrimLeft(arg, "-")]
			if !ok {
				formatError(fmt.Errorf("unknown option: '%s'", arg))
				return
			}

			option()
		}

		// Run fizz file
		filename := args[1]
		if !strings.HasSuffix(filename, ".fizz") {
			filename += ".fizz"
		}

		if err := runFile(filename); err != nil {
			formatError(err)
		}

		return
	}

	// Run terminal mode
	runTerminal()
}

// Prints errors with red color to terminal
func formatError(err error) {
	ct.Foreground(ct.Red, true)
	fmt.Println(err.Error())
	ct.ResetColor()
}

// Leaves the interpreter running as the user inputs code to the terminal.
// Prints out errors but does not terminate until ^C or 'exit'.
func runTerminal() {
	scanner := bufio.NewScanner(os.Stdin)
	totalString := ""
	numBlocks := 0
	indent := "    "

	fmt.Println("type 'exit' to terminate session")
	for {
		fmt.Print("::: " + strings.Repeat(indent, numBlocks))
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Continue with indent after braces
		numBlocks += strings.Count(input, "{") - strings.Count(input, "}")
		totalString += input + "\n" // Better error handling

		if numBlocks <= 0 {
			err := Interperate(totalString)
			if err != nil {
				formatError(err)
			}

			totalString = ""
			numBlocks = 0
		}
	}

	fmt.Println("session ended")
}

// Interperates code found in file specified in commandline arguments
func runFile(filename string) (err error) {
	if !strings.HasSuffix(filename, ".fizz") {
		return ErrNonFizzFile
	}

	if file, err := os.Open(filename); err == nil {
		var buf bytes.Buffer
		bufio.NewReader(file).WriteTo(&buf)
		return Interperate(buf.String())
	}

	// Assumes path error
	return fmt.Errorf(ErrFileNotFound.Error(), filename)
}
