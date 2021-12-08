package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/interp"
	"github.com/jesperkha/Fizz/util"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

// Todo: make simple cmd flag parser
// Have full length versions of all flags. One letter flags for most common
// Flags:
// 		--json print data as json
// 		--no-color prints errors without color
var cmdOptions = map[string]func(){
	// Todo: desc of what program does. a couple run examples. desc of flags. link to github
	"help": func() {
		fmt.Println("use: fizz [--option] [-flag] [filename]")
	},
	"version": func() {
		fmt.Printf("Fizz %s\n", VERSION)
	},
	"flags": func() {
		
	},
}

// Todo: set non 0 exit code for program failure. map to universal codes
func RunInterpeter(args []string) {
	if len(args) == 2 {
		arg := args[1]

		if strings.HasPrefix(arg, "--") {
			if option, ok := cmdOptions[arg[2:]]; ok {
				option()
			} else {
				util.PrintError(fmt.Errorf("unknown option: '%s', use --options to see a full list of options", arg))
			}

			return
		}

		filename := args[1]
		if !strings.HasSuffix(filename, ".fizz") {
			filename += ".fizz"
		}

		if err := runFile(filename); err != nil {
			util.PrintError(err)
		}

		return
	}

	runTerminal()
}

// Leaves the interpreter running as the user inputs code to the terminal.
// Prints out errors but does not terminate until ^C or 'exit'.
func runTerminal() {
	fmt.Println("type 'exit' to terminate session")
	scanner := bufio.NewScanner(os.Stdin)

	totalString := ""
	numBlocks := 0
	line := 1
	space := " "

	for {
		fmt.Printf("%d%s : %s", line, space, strings.Repeat("    ", numBlocks))
		scanner.Scan()
		input := scanner.Text()

		if input == "exit" {
			break
		}

		// Continue with indent after braces
		numBlocks += strings.Count(input, "{") - strings.Count(input, "}")
		totalString += input + "\n" // Better error handling

		if numBlocks <= 0 {
			if _, err := interp.Interperate(totalString); err != nil {
				util.PrintError(err)
				line--
			}

			totalString = ""
			numBlocks = 0
		}

		line++
		if line == 10 {
			space = ""
		}
	}

	fmt.Println("session ended")
}

// Interperates code found in file specified in commandline arguments
func runFile(filename string) (err error) {
	if !strings.HasSuffix(filename, ".fizz") {
		return ErrNonFizzFile
	}

	if byt, err := os.ReadFile(filename); err == nil {
		_, err = interp.Interperate(string(byt))
		return err
	}

	// Assumes path error
	return fmt.Errorf(ErrFileNotFound.Error(), filename)
}
