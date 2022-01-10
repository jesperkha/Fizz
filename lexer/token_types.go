package lexer

const (
	// Ordered by precedence lo -> hi
	NOT_TOKEN = iota

	// Expression types
	AND
	OR

	EQUAL_EQUAL
	NOT_EQUAL
	IN

	GREATER
	LESS
	GREATER_EQUAL
	LESS_EQUAL

	PLUS
	MINUS
	TYPE
	NOT

	STAR
	SLASH
	MODULO
	HAT

	STRING
	NUMBER
	TRUE
	FALSE
	NIL
	IDENTIFIER

	// Not valid expression types
	LEFT_PAREN
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_SQUARE
	RIGHT_SQUARE
	COMMA
	DOT
	SEMICOLON
	COMMENT
	EQUAL
	DEF_EQUAL
	PLUS_EQUAL
	MINUS_EQUAL
	MULT_EQUAL
	DIV_EQUAL

	// Statement tokens
	FUNC
	DEFINE
	ENUM
	IF
	ELSE
	RANGE
	PRINT
	EXIT
	ERROR
	VAR
	WHILE
	BREAK
	SKIP
	IMPORT
	INCLUDE
	REPEAT
	RETURN

	WHITESPACE
	NEWLINE
	EOF
)

var tokenLookup = map[rune]int{
	'\t': WHITESPACE,
	'\r': WHITESPACE,
	'\n': NEWLINE,
	' ':  WHITESPACE,

	')': RIGHT_PAREN,
	'(': LEFT_PAREN,
	'{': LEFT_BRACE,
	'}': RIGHT_BRACE,
	'[': LEFT_SQUARE,
	']': RIGHT_SQUARE,

	'*': STAR,
	'/': SLASH,
	'-': MINUS,
	'+': PLUS,
	'^': HAT,
	'%': MODULO,

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
	"+=": PLUS_EQUAL,
	"-=": MINUS_EQUAL,
	"*=": MULT_EQUAL,
	"/=": DIV_EQUAL,
	":=": DEF_EQUAL,
}

var keyWordLookup = map[string]int{
	"if":      IF,
	"else":    ELSE,
	"nil":     NIL,
	"print":   PRINT,
	"return":  RETURN,
	"exit":    EXIT,
	"define":  DEFINE,
	"var":     VAR,
	"true":    TRUE,
	"false":   FALSE,
	"while":   WHILE,
	"func":    FUNC,
	"type":    TYPE,
	"break":   BREAK,
	"import":  IMPORT,
	"repeat":  REPEAT,
	"skip":    SKIP,
	"include": INCLUDE,
	"error":   ERROR,
	"in":      IN,
	"enum":    ENUM,
	"range":   RANGE,
}
