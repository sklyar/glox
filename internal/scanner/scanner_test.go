package scanner

import (
	"github.com/sklyar/glox/internal/token"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScanSingleCharTokens(t *testing.T) {
	t.Parallel()

	input := "(){};,+-*/"
	expected := []token.Type{
		token.LPAREN, token.RPAREN,
		token.LBRACE, token.RBRACE,
		token.SEMICOLON, token.COMMA,
		token.ADD, token.SUB,
		token.MUL, token.QUO,
		token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	for i, tkn := range tokens {
		require.Equalf(t, expected[i], tkn.Type, "Unexpected token at index %d: %+v", i, tkn)
	}
}

func TestScanMultiCharTokens(t *testing.T) {
	t.Parallel()

	input := "!= == <= >="
	expected := []token.Type{
		token.NEQ, token.EQL,
		token.LEQ, token.GEQ,
		token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	for i, tkn := range tokens {
		require.Equalf(t, expected[i], tkn.Type, "Unexpected token at index %d: %+v", i, tkn)
	}
}

func TestScanStringLiteral(t *testing.T) {
	t.Parallel()

	input := `"hello world"`
	expected := []token.Type{
		token.STRING, token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	require.Equal(t, expected[0], tokens[0].Type)
	require.Equal(t, "hello world", tokens[0].Literal)
}

func TestScanNumberLiteral(t *testing.T) {
	t.Parallel()

	input := "123.45"
	expected := []token.Type{
		token.NUMBER, token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	require.Equal(t, expected[0], tokens[0].Type)
	require.Equal(t, "123.45", string(tokens[0].Literal.([]byte)))
}

func TestScanIdentifiers(t *testing.T) {
	t.Parallel()

	input := "var foo = 42;"
	expected := []token.Type{
		token.VAR, token.IDENT,
		token.ASSIGN, token.NUMBER,
		token.SEMICOLON, token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	require.Equal(t, expected[0], tokens[0].Type)
	require.Equal(t, expected[1], tokens[1].Type)
	require.Equal(t, "foo", tokens[1].Lexeme)
}

func TestScanIllegalToken(t *testing.T) {
	t.Parallel()

	input := "@"
	expected := []token.Type{
		token.ILLEGAL, token.EOF,
	}

	tokens, err := NewScanner([]byte(input), noErrorHandler(t)).ScanTokens()
	require.NoError(t, err)

	require.Len(t, tokens, len(expected))
	require.Equal(t, expected[0], tokens[0].Type)
}

func noErrorHandler(t *testing.T) ErrorHandler {
	t.Helper()

	return func(offset, line int, msg string) {
		t.Errorf("Unexpected error: %s. Line: %d. Offset: %d", msg, line, offset)
	}
}
