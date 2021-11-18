package run

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
		numBlocks += strings.Count(input, "{") - strings.Count(input, "}")
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