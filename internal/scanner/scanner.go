package scanner

import (
	"github.com/sklyar/glox/internal/token"
	"unicode/utf8"
)

const (
	bom = 0xFEFF // byte order mark, only permitted as very first character
	eof = -1     // end of file
)

type ErrorHandler func(offset, line int, msg string)

type Scanner struct {
	src        []byte
	errHandler ErrorHandler

	cur            rune // current character.
	curOffsetStart int  // current character position start offset.
	curOffsetEnd   int  // current character position end offset.
	line           int  // current line number.
}

func NewScanner(b []byte, errHandler ErrorHandler) *Scanner {
	s := &Scanner{
		src:        b,
		errHandler: errHandler,
		line:       1,
	}
	if s.peek() == bom {
		s.next()
	}
	return s
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	tokens := make([]token.Token, 0)
	for !s.isAtEnd() {
		s.curOffsetStart = s.curOffsetEnd

		tkn := s.Scan()
		tokens = append(tokens, tkn)
		if tkn.Type == token.EOF {
			return tokens, nil
		}
	}

	tokens = append(tokens, token.NewToken(token.EOF, "", nil, s.line))
	return tokens, nil
}

func (s *Scanner) Scan() token.Token {
scanAgain:
	s.next()

	switch s.cur {
	case eof:
		return s.newToken(token.EOF)
	case ' ', '\t', '\n', '\r':
		goto scanAgain
	case '(':
		return s.newToken(token.LPAREN)
	case ')':
		return s.newToken(token.RPAREN)
	case '{':
		return s.newToken(token.LBRACE)
	case '}':
		return s.newToken(token.RBRACE)
	case ',':
		return s.newToken(token.COMMA)
	case '.':
		return s.newToken(token.DOT)
	case '-':
		return s.newToken(token.SUB)
	case '+':
		return s.newToken(token.ADD)
	case ';':
		return s.newToken(token.SEMICOLON)
	case '*':
		return s.newToken(token.MUL)
	case '!':
		if s.match('=') {
			return s.newToken(token.NEQ)
		} else {
			return s.newToken(token.NOT)
		}
	case '=':
		if s.match('=') {
			return s.newToken(token.EQL)
		} else {
			return s.newToken(token.ASSIGN)
		}
	case '<':
		if s.match('=') {
			return s.newToken(token.LEQ)
		} else {
			return s.newToken(token.LSS)
		}
	case '>':
		if s.match('=') {
			return s.newToken(token.GEQ)
		} else {
			return s.newToken(token.GTR)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.next()
			}
			// Skip comments.
			goto scanAgain
		} else {
			return s.newToken(token.QUO)
		}
	case '"':
		return s.scanString()
	default:
		if isDecimal(s.cur) {
			return s.scanNumber()
		}
		if isLetter(s.cur) {
			return s.scanIdentifier()
		}
	}

	return s.newToken(token.ILLEGAL)
}

func (s *Scanner) isAtEnd() bool {
	return s.curOffsetEnd >= len(s.src)
}

func (s *Scanner) newToken(tokenType token.Type) token.Token {
	return s.newTokenWithLiteral(tokenType, nil)
}

func (s *Scanner) newTokenWithLiteral(tokenType token.Type, literal any) token.Token {
	lexeme := string(s.src[s.curOffsetStart:s.curOffsetEnd])
	tkn := token.NewToken(tokenType, lexeme, literal, s.line)
	return tkn
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if rune(s.src[s.curOffsetEnd]) != expected {
		return false
	}

	s.curOffsetEnd++
	return true
}

func (s *Scanner) scanString() token.Token {
	startOffset := s.curOffsetStart

	for {
		s.next()
		if s.cur == '\n' || s.cur < 0 {
			s.error(startOffset, "unterminated string")
			break
		}
		if s.cur == '"' {
			break
		}
	}

	s.curOffsetStart = startOffset
	lit := string(s.src[startOffset+1 : s.curOffsetEnd-1])
	return s.newTokenWithLiteral(token.STRING, lit)
}

func (s *Scanner) scanNumber() token.Token {
	startOffset := s.curOffsetStart
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

	lit := s.src[startOffset:s.curOffsetEnd]
	return s.newTokenWithLiteral(token.NUMBER, lit)
}

func (s *Scanner) scanIdentifier() token.Token {
	startOffset := s.curOffsetStart

	for {
		ch := s.peek()
		if !isLetter(ch) && !isDecimal(ch) {
			break
		}
		if ch == 0 || ch == eof {
			break
		}
		s.next()
		if s.cur == eof {
			break
		}
	}

	s.curOffsetStart = startOffset
	lit := string(s.src[s.curOffsetStart:s.curOffsetEnd])
	tokenType, ok := token.LookupIdent(lit)
	if !ok {
		tokenType = token.IDENT
	}
	return s.newToken(tokenType)
}

func (s *Scanner) next() {
	if s.isAtEnd() {
		s.cur = eof
		return
	}

	s.curOffsetStart = s.curOffsetEnd

	r, w := rune(s.src[s.curOffsetStart]), 1
	if r == 0 {
		s.error(s.curOffsetStart, "NUL byte in source")
		return
	}
	if r >= utf8.RuneSelf {
		// not ASCII
		r, w = utf8.DecodeRune(s.src[s.curOffsetStart:])
		if r == utf8.RuneError && w == 1 {
			s.error(s.curOffsetStart, "invalid UTF-8 encoding")
			return
		}
	}

	s.cur = r
	s.curOffsetEnd += w

	return
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return eof
	}

	r := rune(s.src[s.curOffsetEnd])
	if r == 0 {
		return 0
	}

	if r >= utf8.RuneSelf {
		// not ASCII
		var w int
		r, w = utf8.DecodeRune(s.src[s.curOffsetEnd:])
		if r == utf8.RuneError && w == 1 {
			return 0
		}
	}

	return r
}

func (s *Scanner) nextPeek() rune {
	if s.curOffsetEnd+1 >= len(s.src) {
		return eof
	}

	r := rune(s.src[s.curOffsetEnd+1])
	if r == 0 {
		return 0
	}

	if r >= utf8.RuneSelf {
		// not ASCII
		var w int
		r, w = utf8.DecodeRune(s.src[s.curOffsetEnd+1:])
		if r == utf8.RuneError && w == 1 {
			return 0
		}
	}

	return r
}

func (s *Scanner) error(offset int, msg string) {
	s.errHandler(offset, s.line, msg)
}

func isDecimal(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}
