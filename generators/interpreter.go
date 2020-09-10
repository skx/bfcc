package generators

import (
	"fmt"
	"os"

	"github.com/skx/bfcc/lexer"
)

// Interpreter is a generator which will actually interpret, or execute,
// the supplied program.  It is not a compiler, and it produces no
// output (other than that printed by the program itself).
//
// As our host application is designed to be a compiler this is a little
// atypical, however it demonstrates the speedup possible by compilation.
//
// In my rough tests this interpreter is 200 times slower executing our
// mandelbrot example.
type Interpreter struct {

	//
	// Our state
	//

	// The lexed tokens of our input program.
	tokens []*lexer.Token

	// The offset within the program which we're executing.
	offset int

	//
	// BF virtual machine
	//

	// The index pointer.
	ptr int

	// The memory.
	memory [3000]int
}

// Generate takes the specified input-program, and executes it.
func (i *Interpreter) Generate(input string, output string) error {

	// Create a lexer
	lex := lexer.New(input)

	// Store the programs' lexed tokens
	tok := lex.Next()
	for tok.Type != lexer.EOF {
		i.tokens = append(i.tokens, tok)
		tok = lex.Next()
	}

	// Setup our defaults
	i.ptr = 0
	i.offset = 0

	//
	// Repeatedly evaluate a single instruction, until
	// we've exhausted our program.
	//
	for i.offset < len(i.tokens) {
		err := i.evaluate()
		if err != nil {
			return err
		}
	}

	return nil
}

// evaluate executes the current BF instruction.
func (i *Interpreter) evaluate() error {

	//
	// Get the token we're executing.
	//
	tok := i.tokens[i.offset]

	//
	// Execute it.
	//
	switch tok.Type {

	case lexer.INC_PTR:
		i.ptr += tok.Repeat

	case lexer.DEC_PTR:
		i.ptr -= tok.Repeat

	case lexer.INC_CELL:
		i.memory[i.ptr] += tok.Repeat

	case lexer.DEC_CELL:
		i.memory[i.ptr] -= tok.Repeat

	case lexer.LOOP_OPEN:
		// early termination
		if i.memory[i.ptr] != 0 {
			i.offset++
			return nil
		}

		// Otherwise we need to the end of the loop and
		// jump to it
		depth := 1
		for depth != 0 {
			i.offset++
			switch i.tokens[i.offset].Type {
			case lexer.LOOP_OPEN:
				depth++
			case lexer.LOOP_CLOSE:
				depth--
			}
		}
		return nil

	case lexer.LOOP_CLOSE:

		// early termination
		if i.memory[i.ptr] == 0 {
			i.offset++
			return nil
		}

		depth := 1
		for depth != 0 {
			i.offset--
			switch i.tokens[i.offset].Type {
			case lexer.LOOP_CLOSE:
				depth++
			case lexer.LOOP_OPEN:
				depth--
			}
		}
		return nil

	case lexer.INPUT:
		buf := make([]byte, 1)
		l, err := os.Stdin.Read(buf)
		if err != nil {
			return err
		}
		if l != 1 {
			return fmt.Errorf("read %d bytes of input, not 1", l)
		}
		i.memory[i.ptr] = int(buf[0])

	case lexer.OUTPUT:
		fmt.Printf("%c", rune(i.memory[i.ptr]))

	}

	// next instruction will be executed next time.
	i.offset++
	return nil
}

// Register our back-end
func init() {
	Register("interpreter", func() Generator {
		return &Interpreter{}
	})
}
