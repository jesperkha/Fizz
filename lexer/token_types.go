package lexer

const (
	// Ordered by precedence lo -> hi
	NOT_TOKEN = iota
	AND
	OR

	EQUAL_EQUAL
	NOT_EQUAL

	GREATER
	LESS
	GREATER_EQUAL
	LESS_EQUAL

	PLUS
	MINUS
	TYPE

	STAR
	SLASH
	HAT

	NOT

	IDENTIFIER
	STRING
	NUMBER

	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	SEMICOLON
	COMMENT
	EQUAL
	FUNC
	CLASS
	RETURN
	IF
	ELSE
	FOR
	NIL
	PRINT
	SUPER
	THIS
	TRUE
	FALSE
	VAR
	WHILE
	BREAK
	WHITESPACE
	NEWLINE
	EOF
)

var tokenLookup = map[rune]int{
	'\n': NEWLINE,
	'\t': WHITESPACE,
	'\r': WHITESPACE,
	' ':  WHITESPACE,

	')': RIGHT_PAREN,
	'(': LEFT_PAREN,
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,

	'*': STAR,
	'/': SLASH,
	'-': MINUS,
	'+': PLUS,
	'^': HAT,

	';': SEMICOLON,
	',': COMMA,
	'.': DOT,
	'#': COMMENT,
	'"': STRING,

	'=': EQUAL,
	'!': NOT,
	'>': GREATER,
	'<': LESS,
	'&': AND,
	':': OR,
}

var doubleTokenLookup = map[string]int{
	"==": EQUAL_EQUAL,
	"!=": NOT_EQUAL,
	">=": GREATER_EQUAL,
	"<=": LESS_EQUAL,
}

var keyWordLookup = map[string]int{
	"if":     IF,
	"else":   ELSE,
	"nil":    NIL,
	"print":  PRINT,
	"return": RETURN,
	"class":  CLASS,
	"for":    FOR,
	"super":  SUPER,
	"this":   THIS,
	"var":    VAR,
	"true":   TRUE,
	"false":  FALSE,
	"while":  WHILE,
	"func":   FUNC,
	"type":   TYPE,
	"break":  BREAK,
}