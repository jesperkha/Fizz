package lexer

// Tokenizes text input into slice. Errors can be raised
// for invalid tokens, identifiers, or unlosed strings.

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnexpectedToken    = errors.New("unexpeted token: '%s', line: %d")
	ErrUnterminatedString = errors.New("unterminated string, line %d")
	ErrInvalidSyntax      = errors.New("invalid syntax '%s', line %d")
)

type Token struct {
	Type    int
	Lexeme  string
	Literal interface{}
	Line    int
}

func GetTokens(input string) (tokens []Token, err error) {
	currentIdx := 0
	currentLine := 1 // Start at line 1 for editors
	alphaNumRegex := regexp.MustCompile("^[a-zA-Z_0-9.]*$")
	variableRegex := regexp.MustCompile("^[a-zA-Z_][a-zA-Z_0-9]*$")

	for currentIdx < len(input) {
		startIndex := currentIdx
		char := input[currentIdx]
		nextChar, _ := getNextCharacter(input, currentIdx)

		tokenType, isSymbol := tokenLookup[rune(char)]
		token := Token{Type: tokenType, Lexeme: string(char), Line: currentLine}

		if isSymbol {
			// Token exception cases
			switch tokenType {
			case NEWLINE:
				currentLine++
				currentIdx++
				tokens = append(tokens, Token{Type: NEWLINE})
				continue
			case WHITESPACE:
				currentIdx++
				continue
			case COMMENT:
				seekCharacter(input, &currentIdx, '\n')
				continue
			}

			// Check for double symbol (!=, >= etc)
			if nextType, ok := tokenLookup[nextChar]; ok && nextType == EQUAL {
				jointSymbol := strings.Join([]string{string(char), string(nextChar)}, "")
				token.Lexeme = jointSymbol
				token.Type = doubleTokenLookup[jointSymbol]
				currentIdx++ // Skip next char
			}

			// Seek closing string
			if tokenType == STRING {
				if seekCharacter(input, &currentIdx, '"') {
					return tokens, fmt.Errorf(ErrUnterminatedString.Error(), currentLine)
				}

				str := intervalToString(input, startIndex, currentIdx)
				token.Lexeme = str
				token.Literal = str[1 : len(str)-1]
			}

			tokens = append(tokens, token)
			currentIdx++
			continue
		}

		// Not alpha numeric (a-z 0-9 _.)
		if !alphaNumRegex.MatchString(string(char)) {
			return tokens, fmt.Errorf(ErrUnexpectedToken.Error(), string(char), currentLine)
		}

		// Char is not a symbol and is the start of an identifier, keyword, or number
		seekFunc(input, &currentIdx, func(c rune) bool {
			return !alphaNumRegex.MatchString(string(c))
		})

		identifier := intervalToString(input, startIndex, currentIdx)
		number, err := strconv.ParseFloat(identifier, 64)
		token.Lexeme = identifier

		isNumber := err == nil
		isAlphaNum := variableRegex.MatchString(identifier)

		splitDot := strings.Split(identifier, ".")
		isGetter := !isNumber && len(splitDot) > 1

		invalidSyntax := fmt.Errorf(ErrInvalidSyntax.Error(), identifier, currentLine)
		if !isNumber && !isAlphaNum && !isGetter {
			return tokens, invalidSyntax
		}

		if isGetter {
			dot, err := GetTokens(".")
			if err != nil {
				return tokens, err
			}

			ts := []Token{}
			for _, ident := range splitDot {
				t, err := GetTokens(ident)
				if err != nil {
					return tokens, err
				}

				if len(t) == 0 {
					return tokens, invalidSyntax
				}

				t[0].Line = token.Line
				ts = append(ts, dot...)
				ts = append(ts, t...)
			}

			// Shift 1 to skip first dot
			tokens = append(tokens, ts[1:]...)
			currentIdx++
			continue
		}

		if isNumber {
			token.Literal = number
			token.Type = NUMBER
		}

		if isAlphaNum {
			token.Type = IDENTIFIER
			if keywordType, isKeyword := keyWordLookup[identifier]; isKeyword {
				// Set literal values for keyword types
				switch keywordType {
				case FALSE:
					token.Literal = false
				case TRUE:
					token.Literal = true
				case NIL:
					token.Literal = nil
				}

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
	if curIdx < len(input)-1 {
		return rune(input[curIdx+1]), false
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
