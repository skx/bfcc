package lexer

import (
	"testing"
)

// TestLexer performs a trivial test of the lexer
func TestLexer(t *testing.T) {

	tests := []struct {
		expectedType  string
		expectedCount int
	}{
		{PLUS, 1},
		{MINUS, 1},
		{LESS, 5},
		{GREATER, 5},
		{LOOPOPEN, 1},
		{LOOPCLOSE, 1},
		{OUTPUT, 1},
		{INPUT, 1},
		{EOF, 1},
	}

	l := NewLexer("+-<<<<<\n>>>>>[].,")

	for i, tt := range tests {
		tok := l.Next()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Repeat != tt.expectedCount {
			t.Fatalf("tests[%d] - count wrong, expected=%d, got=%d", i, tt.expectedCount, tok.Repeat)
		}
	}
}

// TestAdjacent is designed to ensure we count adjacent runs of characters
// even when newlines are in the way.
func TestAdjacent(t *testing.T) {

	tests := []struct {
		expectedType  string
		expectedCount int
	}{
		{PLUS, 5},
		{MINUS, 5},
		{EOF, 1},
	}

	l := NewLexer("+\n+\n+\n+\n+- - - - -")

	for i, tt := range tests {
		tok := l.Next()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong, expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Repeat != tt.expectedCount {
			t.Fatalf("tests[%d] - count wrong, expected=%d, got=%d", i, tt.expectedCount, tok.Repeat)
		}
	}
}
