package main

import (
	"testing"
)

func TestLexer(t *testing.T) {

	tests := []struct {
		expectedType  string
		expectedCount int
	}{
		{PLUS, 1},
		{MINUS, 1},
		{LESS, 5},
		{GREATER, 5},
		{LOOP_OPEN, 1},
		{LOOP_CLOSE, 1},
		{OUTPUT, 1},
		{INPUT, 1},
		{EOF, 1},
	}

	l := NewLexer("+-<<<<<>>>>>[].,")

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
