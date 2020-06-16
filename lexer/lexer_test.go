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
		{INC_CELL, 1},
		{DEC_CELL, 1},
		{DEC_PTR, 5},
		{INC_PTR, 5},
		{LOOP_OPEN, 1},
		{LOOP_CLOSE, 1},
		{OUTPUT, 1},
		{INPUT, 1},
		{EOF, 1},
	}

	l := New("+-<<<<<#\n`>>>>>[].,")

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
		{INC_CELL, 5},
		{DEC_CELL, 5},
		{EOF, 1},
	}

	l := New("+\n+\n+\n+\n+- - - - -")

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
