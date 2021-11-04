package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

var cmdOptions = map[string]func() {
	"help": func() {
		fmt.Println("use: fizz [filename.fizz | --option]")
	},
}

// Interperates code found in file specified in commandline arguments
func RunFile(filename string) (err error) {
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

// Leaves the interpreter running as the user inputs code to the terminal.
// Prints out errors but does not terminate until ^C or 'exit'.
func RunTerminal() {
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
		numBlocks += strings.Count(input, "{")
		numBlocks -= strings.Count(input, "}")
		totalString += input + "\n" // Better error handling

		if numBlocks <= 0 {
			err := Interperate(totalString)
			if err != nil {
				fmt.Println(err.Error() + "\n")
			}

			totalString = ""
			numBlocks = 0
		}
	}

	fmt.Println("session ended")
}

func main() {
	if len(os.Args) == 2 {
		arg := os.Args[1]

		// Run option commands
		if strings.HasPrefix(arg, "--") {
			if opt, ok := cmdOptions[strings.TrimLeft(arg, "-")]; ok {
				opt()
				return
			}

			fmt.Printf("unknown option: '%s'\n", arg)
			return
		}

		// Run fizz file
		filename := os.Args[1]
		if !strings.HasSuffix(filename, ".fizz") {
			filename += ".fizz"
		}

		if err := RunFile(filename); err != nil {
			fmt.Println(err)
		}

		return
	}

	// Run terminal mode
	RunTerminal()
}
