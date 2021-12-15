package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jesperkha/Fizz/interp"
	"github.com/jesperkha/Fizz/stmt"
	"github.com/jesperkha/Fizz/term"
	"github.com/jesperkha/Fizz/util"
)

var parser = term.NewFlagParser(
	[]string{"v"},
	[]string{"version", "help"},
)

func init() {
	parser.Assign("--version", func ()  {
		fmt.Printf("Fizz %s\n", VERSION)
		os.Exit(0)
	})

	parser.Assign("--help", func ()  {
		fmt.Println(term.HELP)
		os.Exit(0)
	})
}

func RunInterpeter(args []string) {
	filename, err := parser.Parse()
	if err != nil {
		util.ErrorAndExit(err)
	}

	if filename != "" {
		// Goto directory of file specified
		split := strings.Split(filename, "/")
		path := strings.Join(split[:len(split)-1], "/")
		name := split[len(split)-1]
		os.Chdir(path)
		if _, err = interp.RunFile(name); err != nil && err != stmt.ErrProgramExit {
			util.ErrorAndExit(err)
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
			if _, err := interp.Interperate("", totalString); err != nil {
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