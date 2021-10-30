package lexer

import (
	"errors"
	"fmt"
	"strings"
)

// Tokenize input text

var (
	ErrUnexpectedToken    = errors.New("unexpeted token: '%s', line: %d")
	ErrInvalidSyntax	  = errors.New("invalid syntax: '%s', line: %d")
	ErrUnterminatedString = errors.New("unterminated string, line %d")
	ErrSyntaxError		  = errors.New("invalid syntax, line %d")
)

type Token struct {
	Type 	int
	Lexeme 	string
	Literal interface{}
	Line	int
}

func GetTokens(input string) (tokens []Token, err error) {
	currentIdx  := 0
	currentLine := 0

	for currentIdx < len(input) {
		startIndex := currentIdx
		char := input[currentIdx]
		nextChar, _ := getNextCharacter(input, currentIdx)

		tokenType, isSymbol := tokenLookup[rune(char)]
		lexeme := string(char)
		
		if isSymbol {
			// Check for double symbol (!=, >= etc)
			if nextType, ok := tokenLookup[nextChar]; ok && nextType == EQUAL {
				jointSymbol := strings.Join([]string{string(char), string(nextChar)}, "")
				lexeme = jointSymbol
				tokenType = doubleTokenLookup[jointSymbol]
				currentIdx++ // Skip next char
			}
			
			// Goto new line
			if tokenType == NEWLINE {
				currentLine++
			}

			// Skip token generation
			if tokenType == WHITESPACE {
				currentIdx++
				continue
			}

			// Seek closing string
			if tokenType == STRING {
				eof := seekCharacter(input, &currentIdx, '"')
				if eof {
					return tokens, fmt.Errorf(ErrUnterminatedString.Error(), currentLine)
				}

				lexeme = intervalToString(input, startIndex, currentIdx)
			}
		}

		// For loop checks if alphanumeric and also checks + assigns keyword token

		tokens = append(tokens, Token{Type: tokenType, Lexeme: lexeme, Line: currentLine})
		currentIdx++
	}

	// tokens = append(tokens, Token{Type: EOF})
	return tokens, err
}

// Returns the next character in the input without consuming it
func getNextCharacter(input string, curIdx int) (nextChar rune, eof bool) {
	if curIdx < len(input) - 1 {
		return rune(input[curIdx + 1]), false
	}

	return nextChar, true
}

// Consumes characters until it reaches the target. If the it reaches eof it is returned true. Consumes target.
func seekCharacter(input string, curIdx *int, target rune) (eof bool) {
	for {
		if nextChar, isEOF := getNextCharacter(input, *curIdx); !isEOF {
			if nextChar == target {
				*curIdx++ // Skip final char
				return false
			}

			*curIdx++
			continue
		}

		return true
	}
}

// Takes an index interval in the input and returns the string.
func intervalToString(input string, startIdx int, endIdx int) string {
	result := ""
	for i := startIdx; i <= endIdx; i++ {
		result += string(input[i])
	}
	
	return result
}
