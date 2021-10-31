package lexer

const (
	// Ordered by precedence lo -> hi
	NOT_TOKEN = iota
	EQUAL_EQUAL
	NOT_EQUAL

	GREATER
	LESS
	GREATER_EQUAL
	LESS_EQUAL

	PLUS
	MINUS

	STAR
	SLASH

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
	AND
	OR
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
	'}': LEFT_BRACE,
	'{': RIGHT_BRACE,

	'*': STAR,
	'/': SLASH,
	'-': MINUS,
	'+': PLUS,

	';': SEMICOLON,
	',': COMMA,
	'.': DOT,
	'#': COMMENT,
	'"': STRING,

	'=': EQUAL,
	'!': NOT,
	'>': GREATER,
	'<': LESS,
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
	"or":     OR,
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
}