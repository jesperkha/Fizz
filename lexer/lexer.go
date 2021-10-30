package lexer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Tokenize input text

var (
	ErrUnexpectedToken    = errors.New("unexpeted token: '%s', line: %d")
	ErrUnterminatedString = errors.New("unterminated string, line %d")
	ErrInvalidIdentifier  = errors.New("invalid identifier '%s', line %d")
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
	alphaNumRegex := regexp.MustCompile("^[a-zA-Z_0-9]*$")
	variableRegex := regexp.MustCompile("^[a-zA-Z_][a-zA-Z_0-9]*$")

	for currentIdx < len(input) {
		startIndex := currentIdx
		char := input[currentIdx]
		nextChar, _ := getNextCharacter(input, currentIdx)

		tokenType, isSymbol := tokenLookup[rune(char)]
		token := Token{Type: tokenType, Lexeme: string(char)}
		
		if isSymbol {
			// Check for double symbol (!=, >= etc)
			if nextType, ok := tokenLookup[nextChar]; ok && nextType == EQUAL {
				jointSymbol := strings.Join([]string{string(char), string(nextChar)}, "")
				token.Lexeme = jointSymbol
				tokenType = doubleTokenLookup[jointSymbol]
				currentIdx++ // Skip next char
			}
			
			// Goto new line
			if tokenType == NEWLINE {
				currentLine++
				currentIdx++
				continue
			}

			// Skip token generation
			if tokenType == WHITESPACE {
				currentIdx++
				continue
			}

			// Skip comment
			if tokenType == COMMENT {
				seekCharacter(input, &currentIdx, '\n')
				currentLine++
				continue
			}

			// Seek closing string
			if tokenType == STRING {
				if seekCharacter(input, &currentIdx, '"') {
					return tokens, fmt.Errorf(ErrUnterminatedString.Error(), currentLine)
				}

				token.Lexeme = intervalToString(input, startIndex, currentIdx)
			}

			tokens = append(tokens, token)
			currentIdx++
			continue
		}

		// Not alpha numeric (a-z 0-9 _)
		if !alphaNumRegex.MatchString(string(char)) {
			return tokens, fmt.Errorf(ErrUnexpectedToken.Error(), string(char), currentLine)
		}

		// Char is not a symbol and is the start of an identifier, keyword, or number
		seekFunc(input, &currentIdx, func(c rune) bool {
			return !alphaNumRegex.MatchString(string(c))
		})
		
		identifier := intervalToString(input, startIndex, currentIdx)
		number, err := strconv.Atoi(identifier)
		token.Lexeme = identifier

		isNumber := err == nil
		isAlphaNum := variableRegex.MatchString(identifier)

		if !isNumber && !isAlphaNum {
			return tokens, fmt.Errorf(ErrInvalidIdentifier.Error(), identifier, currentLine)
		}

		if isNumber {
			token.Literal = number
			token.Type = NUMBER
		}
		
		if isAlphaNum {
			token.Type = IDENTIFIER
			if keywordType, isKeyword := keyWordLookup[identifier]; isKeyword {
				token.Type = keywordType
			}
		}

		tokens = append(tokens, token)
		currentIdx++
	}

	return tokens, err
}

// Returns the next character in the input without consuming it
func getNextCharacter(input string, curIdx int) (nextChar rune, eof bool) {
	if curIdx < len(input) - 1 {
		return rune(input[curIdx + 1]), false
	}

	return nextChar, true
}

// Consumes characters until the matchFunc returns true. If eof is reached true is returned.
// Modifies curIdx value. Does not consume final character.
func seekFunc(input string, curIdx *int, matchFunc func(char rune) bool) (eof bool) {
	for {
		if nextChar, isEOF := getNextCharacter(input, *curIdx); !isEOF {
			if matchFunc(nextChar) {
				return false
			}

			*curIdx++
			continue
		}

		return true
	}
}

// Calls seekFunc to match the target character. Consumes final character.
func seekCharacter(input string, curIdx *int, target rune) (eof bool) {
	eof = seekFunc(input, curIdx, func(char rune) bool {
		return char == target
	})

	*curIdx++
	return eof
}

// Takes an index interval in the input and returns the string
func intervalToString(input string, startIdx int, endIdx int) string {
	result := ""
	for i := startIdx; i <= endIdx; i++ {
		result += string(input[i])
	}
	
	return result
}