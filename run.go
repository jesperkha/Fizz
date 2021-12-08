package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/cmd"
	"github.com/jesperkha/Fizz/interp"
	"github.com/jesperkha/Fizz/util"
)

var (
	ErrFileNotFound = errors.New("cannot find file with name: '%s'")
	ErrNonFizzFile  = errors.New("cannot run non-Fizz file")
)

var parser = cmd.NewFlagParser(
	[]string{},
	[]string{"version", "help"},
)

func init() {
	parser.Assign("--version", func ()  {
		fmt.Printf("Fizz %s\n", VERSION)
		os.Exit(0)
	})

	parser.Assign("--help", func ()  {
		fmt.Println("use:\n\tfizz [filename] <flags>\n\tfizz <flags>\n\tfizz")
		os.Exit(0)
	})
}

func RunInterpeter(args []string) {
	filename, err := parser.Parse()
	if err != nil {
		util.ErrorAndExit(err)
	}

	if filename != "" {
		if err = RunFile(filename); err != nil {
			util.ErrorAndExit(fmt.Errorf("%s: %s", filename, err.Error()))
		}

		return
	}

	RunTerminal()
}

// Leaves the interpreter running as the user inputs code to the terminal.
// Prints out errors but does not terminate until ^C or 'exit'.
func RunTerminal() {
	fmt.Println("type 'exit' to terminate session")
	scanner := bufio.NewScanner(os.Stdin)
	numBlocks, line := 0, 1
	totalString, space := "", ""

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

func RunFile(filename string) (err error) {
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
