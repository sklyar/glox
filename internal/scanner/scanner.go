package scanner

import (
	"github.com/sklyar/glox/internal/token"
)

type GeneralError struct {
	Line int
}

func (e *GeneralError) Error() string {
	return "Unterminated string"
}

type UnexpectedCharacterError struct {
	Character rune
	Line      int
}

func (e *UnexpectedCharacterError) Error() string {
	return "Unexpected character"
}

type Scanner struct {
	src    []byte
	tokens []token.Token

	start   int
	current int
	line    int
}

func NewScanner(b []byte) *Scanner {
	return &Scanner{
		src:     b,
		tokens:  nil,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.TypeEOF, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	ch := s.next()
	switch ch {
	case '(':
		s.addToken(token.TypeLeftParen)
	case ')':
		s.addToken(token.TypeRightParen)
	case '{':
		s.addToken(token.TypeLeftBrace)
	case '}':
		s.addToken(token.TypeRightBrace)
	case ',':
		s.addToken(token.TypeComma)
	case '.':
		s.addToken(token.TypeDot)
	case '-':
		s.addToken(token.TypeMinus)
	case '+':
		s.addToken(token.TypePlus)
	case ';':
		s.addToken(token.TypeSemicolon)
	case '*':
		s.addToken(token.TypeStar)
	case '!':
		if s.match('=') {
			s.addToken(token.TypeBangEqual)
		} else {
			s.addToken(token.TypeBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.TypeEqualEqual)
		} else {
			s.addToken(token.TypeEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.TypeLessEqual)
		} else {
			s.addToken(token.TypeLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.TypeGreaterEqual)
		} else {
			s.addToken(token.TypeGreater)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.next()
			}
		} else {
			s.addToken(token.TypeSlash)
		}
	case ' ', '\r', '\t': // Ignore whitespace.
	case '\n':
		s.line++
	case '"':
		if err := s.scanString(); err != nil {
			return err
		}
	default:
		if isDecimal(ch) {
			if err := s.scanNumber(); err != nil {
				return err
			}
			break
		}
		if isLetter(ch) {
			s.scanIdentifier()
		}
		return &UnexpectedCharacterError{Character: ch, Line: s.line}
	}

	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.src)
}

func (s *Scanner) next() rune {
	ch := s.src[s.current]
	s.current++
	return rune(ch)
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}

	return rune(s.src[s.current])
}

func (s *Scanner) addToken(tokenType token.Type) {
	s.addToken2(tokenType, nil)
}

func (s *Scanner) addToken2(tokenType token.Type, literal any) {
	lexeme := string(s.src[s.start:s.current])
	s.tokens = append(s.tokens, token.NewToken(tokenType, lexeme, literal, s.line))
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.src[s.current]) != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) scanString() error {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.next()
	}

	if s.isAtEnd() {
		return &GeneralError{Line: s.line}
	}

	// Consume closing `"`.
	s.next()

	// Trim the surrounding quotes.
	lit := string(s.src[s.start+1 : s.current-1])
	s.addToken2(token.TypeString, lit)

	return nil
}

func (s *Scanner) scanNumber() error {
	for isDecimal(s.peek()) {
		s.next()
	}

	if s.peek() == '.' && isDecimal(s.nextPeek()) {
		// Consume ".".
		s.next()

		for isDecimal(s.peek()) {
			s.next()
		}
	}

	lit := s.src[s.start:s.current]
	s.addToken2(token.TypeNumber, lit)

	return nil
}

func (s *Scanner) scanIdentifier() {
	for ch := s.peek(); isLetter(ch) || isDecimal(ch); {
		s.next()
	}

	lit := string(s.src[s.start:s.current])
	tokenType, ok := token.Lookup(lit)
	if !ok {
		tokenType = token.TypeIdentifier
	}
	s.addToken(tokenType)
}

func (s *Scanner) nextPeek() rune {
	if s.current+1 >= len(s.src) {
		return 0
	}
	return rune(s.src[s.current+1])
}

func isDecimal(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}
