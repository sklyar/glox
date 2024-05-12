package token

type Type uint8

const (
	// Special tokens.
	EOF Type = iota
	ILLEGAL

	// Single-character tokens.
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	COMMA     // ,
	DOT       // .
	SEMICOLON // ;

	ADD // +
	SUB // -
	MUL // *
	QUO // /

	// One or two character tokens.
	EQL    // ==
	LSS    // <
	GTR    // >
	ASSIGN // =
	NOT    // !

	NEQ // !=
	LEQ // <=
	GEQ // >=

	// Literals.
	IDENT  // identifier
	STRING // "string"
	NUMBER // 1234 or 1234.567

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	TRUE
	FUNC
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	VAR
	WHILE
)

var keywords = map[string]Type{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"func":   FUNC,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type    Type
	Lexeme  string
	Literal any // TODO
	Line    int
}

func NewToken(tokenType Type, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func LookupIdent(ident string) (Type, bool) {
	t, ok := keywords[ident]
	return t, ok
}
