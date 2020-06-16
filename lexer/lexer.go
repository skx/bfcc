//lint:file-ignore ST1003 Ignore all-caps because this is fine.

// Package lexer is a simple lexer for processing BrainFuck programs.
//
// A given program will be parsed into a series of tokens, with unrecognized
// input characters skipped over.
//
// As an optimization we collapse multiple adjacent tokens together, and
// return their type as well as a count of how many times the character
// was repeated.
package lexer

import "strings"

// These constants are our token-types
const (
	EOF = "EOF"

	DEC_PTR = "<"
	INC_PTR = ">"

	INC_CELL = "+"
	DEC_CELL = "-"

	OUTPUT = "."
	INPUT  = ","

	LOOP_OPEN  = "["
	LOOP_CLOSE = "]"
)

// Token contains the next token from the input program.
type Token struct {

	// Type contains the token-type (such as "<", "[", etc).
	Type string

	// Repeat contains the number of consecutive appearances we've seen
	// of this token.
	Repeat int
}

// Lexer holds our lexer state.
type Lexer struct {

	// input is the string we're lexing.
	input string

	// position is the current position within the input-string.
	position int

	// simple map of single-character tokens to their type
	known map[string]string

	// simple map which allows us to determine if a token can
	// have repeated occurences collapsed.
	repeat map[string]bool
}

// New creates a new Lexer, which will parse the specified
// input program into a series of tokens.
func New(input string) *Lexer {

	// Create the lexer object.
	l := &Lexer{input: input}

	// Strip newlines/spaces from our iput
	l.input = strings.ReplaceAll(l.input, "\n", "")
	l.input = strings.ReplaceAll(l.input, "\r", "")
	l.input = strings.ReplaceAll(l.input, " ", "")

	// Populate the simple token-types in a map for later use.
	l.known = make(map[string]string)
	l.known["+"] = INC_CELL
	l.known["-"] = DEC_CELL
	l.known[">"] = INC_PTR
	l.known["<"] = DEC_PTR
	l.known[","] = INPUT
	l.known["."] = OUTPUT
	l.known["["] = LOOP_OPEN
	l.known["]"] = LOOP_CLOSE

	// Some characters will have their input collapsed
	// when multiple consecutive occurrences are found.
	l.repeat = make(map[string]bool)
	l.repeat["+"] = true
	l.repeat["-"] = true
	l.repeat[">"] = true
	l.repeat["<"] = true

	return l
}

// Next returns the next token from our input stream.
//
// This is pretty naive lexer because we only have to consider
// single-character tokens.  However we do look for tokens which
// are repeated.
func (l *Lexer) Next() *Token {

	// Loop until we've exhausted our input.
	for l.position < len(l.input) {

		// Get the next character
		char := string(l.input[l.position])

		// Is this a known character/token?
		_, ok := l.known[char]
		if ok {

			//
			// Can this token be repeated?
			//
			// If not just return this single instance.
			//
			repeated := l.repeat[char]
			if !repeated {
				l.position++
				return &Token{Type: char, Repeat: 1}
			}

			//
			// OK we've found a character, such as "<",
			// which can be repeated multiple times.
			//
			// We count how many times that repetition
			// occurs, swallowing that input as we go.
			//

			// Record our starting position
			begin := l.position

			// Loop forward to see how many times the character
			// is repeated.
			for l.position < len(l.input) {

				// If it isn't the same character
				// we're done
				if string(l.input[l.position]) != char {
					break
				}

				// Otherwise keep advancing forward
				l.position++
			}

			// Return the token and the times it was
			// seen in adjacent positions
			count := l.position - begin
			return &Token{Type: char, Repeat: count}
		}

		//
		// Here we're ignoring a token which was unknown.
		//
		l.position++
	}

	//
	// If we got here then we're at/after the end of our input
	// string.  So we just return EOF.
	//
	return &Token{Type: EOF, Repeat: 1}
}
