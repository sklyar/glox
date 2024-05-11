package token

type Type uint8

const (
	// Single-character tokens.
	TypeLeftParen Type = iota + 1
	TypeRightParen
	TypeLeftBrace
	TypeRightBrace
	TypeComma
	TypeDot
	TypeMinus
	TypePlus
	TypeSemicolon
	TypeSlash
	TypeStar

	// One or two character tokens.
	TypeBang
	TypeBangEqual
	TypeEqual
	TypeEqualEqual
	TypeGreater
	TypeGreaterEqual
	TypeLess
	TypeLessEqual

	// Literals.
	TypeIdentifier
	TypeString
	TypeNumber

	// Keywords.
	TypeAnd
	TypeClass
	TypeElse
	TypeFalse
	TypeTrue
	TypeFunc
	TypeFor
	TypeIf
	TypeNil
	TypeOr
	TypePrint
	TypeReturn
	TypeSuper
	TypeThis
	TypeVar
	TypeWhile

	TypeEOF
)

var keywords = map[string]Type{
	"and":    TypeAnd,
	"class":  TypeClass,
	"else":   TypeElse,
	"false":  TypeFalse,
	"for":    TypeFor,
	"func":   TypeFunc,
	"if":     TypeIf,
	"nil":    TypeNil,
	"or":     TypeOr,
	"print":  TypePrint,
	"return": TypeReturn,
	"super":  TypeSuper,
	"this":   TypeThis,
	"true":   TypeTrue,
	"var":    TypeVar,
	"while":  TypeWhile,
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

func Lookup(ident string) (Type, bool) {
	t, ok := keywords[ident]
	return t, ok
}
