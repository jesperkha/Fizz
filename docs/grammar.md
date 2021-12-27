# Fizz grammar

```py
# Fizz program
fizzProgram -> declaration*

# Declarations
declaration -> varDec | funcDec | objDec | statement
varDec      -> identifier ":=" expression ";"
funcDec     -> "func" "(" identifier? ("," identifier)* ")" block
objDec      -> "define" identifier "{"
                    identifier? ("," identifier)* |
                    identifier? (newline identifier)*
                "}"

# Statements
statement   -> exprStmt | printStmt | exitStmt | errorStmt | ifStmt | whileStmt |
               returnStmt | importStmt | includeStmt | assignStmt | block
exprStmt    -> expression ";"
printStmt   -> "print" expression ";"
exitStmt    -> "exit" expression? ";"
errorStmt   -> "error" expression ";"
ifStmt      -> "if" expression block ("else" block)?
whileStmt   -> "while" expression block
returnStmt  -> "return" expression? ";"
importStmt  -> "import" string ";"
includeStmt -> "include" string ";"
assignStmt  -> (getter | identifier) "=" expression ";"
block       -> "{" declaration* "}"

# Expressions
expression -> literal | unary | binary | group | call | array | getter | index
literal    -> "true" | "false" | "nil" | identifier | number | string
unary      -> ("-", "!", "type") expression
binary     -> expression operator expression
group      -> "(" expression ")"
call       -> expression "(" expression? ("," expression)* ")"*
array      -> "[" expression? ("," expression)* "]"
getter     -> expression "." identfier
index      -> array "[" expression "]"

# Operators
operator -> "+" | "-" | "*" | "/" | "^" | "%" | "&" |
            ":" | "==" | "!=" | ">=" | "<=" | "<" | ">"
assignOp -> "=" | ":="
```
