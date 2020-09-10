package generators

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/skx/bfcc/lexer"
)

// GeneratorC is a generator that will produce an C version of the specified
// input-program.
//
// The C-source will then be compiled by gcc.
type GeneratorC struct {
	// input source
	input string

	// file to write to
	output string
}

// generateSource produces a version of the program as C source-file
func (c *GeneratorC) generateSource() error {
	var buff bytes.Buffer
	var programStart = `
extern int putchar(int);
extern char getchar();

char array[30000];

int idx = 0;

int main (int arc, char *argv[]) {
`
	buff.WriteString(programStart)

	//
	// Create a lexer for the input program
	//
	l := lexer.New(c.input)

	//
	// Program consists of all tokens
	//
	program := l.Tokens()

	//
	// We'll process the complete program until
	// we hit an end of file/input
	//
	offset := 0
	for offset < len(program) {

		//
		// The current token
		//
		tok := program[offset]

		//
		// Output different things depending on the token-type
		//
		switch tok.Type {

		case lexer.INC_PTR:
			buff.WriteString(fmt.Sprintf("  idx += %d;\n", tok.Repeat))
		case lexer.DEC_PTR:
			buff.WriteString(fmt.Sprintf("  idx -= %d;\n", tok.Repeat))
		case lexer.INC_CELL:
			buff.WriteString(fmt.Sprintf("  array[idx] += %d;\n", tok.Repeat))
		case lexer.DEC_CELL:
			buff.WriteString(fmt.Sprintf("  array[idx] -= %d;\n", tok.Repeat))
		case lexer.OUTPUT:
			buff.WriteString("  putchar(array[idx]);\n")

		case lexer.INPUT:
			buff.WriteString("  array[idx] = getchar();\n")

		case lexer.LOOP_OPEN:

			//
			// We sneekily optimize "[-]" by converting it
			// into an explicit setting of the cell-content
			// to zero.
			//
			// Since this involves looking at future-tokens
			// we need to make sure we're not at the end of
			// the program.
			//
			if offset+2 < len(program) {

				//
				// Look for the next two tokens "-]", if
				// we find them then we're looking at "[-]"
				// which is something we can optimize.
				//
				if program[offset+1].Type == lexer.DEC_CELL &&
					program[offset+2].Type == lexer.LOOP_CLOSE {
					// register == zero
					buff.WriteString("  array[idx] = 0;\n")

					// 1. Skip this instruction,
					// 2. the next one "-"
					// 3. and the final one "]"
					offset += 3

					// And continue the loop again.
					continue
				}
			}

			buff.WriteString("  while (array[idx]) {\n")

		case lexer.LOOP_CLOSE:
			buff.WriteString("}\n")

		default:
			fmt.Printf("token not handled: %v\n", tok)
			os.Exit(1)
		}

		//
		// Keep processing
		//
		offset++
	}

	// Close the main-function
	buff.WriteString("}\n")

	// Output to a file
	err := ioutil.WriteFile(c.output+".c", buff.Bytes(), 0644)
	return err
}

// compileSource uses gcc to compile the generated source-code
func (c *GeneratorC) compileSource() error {

	gcc := exec.Command(
		"gcc",
		"-static",
		"-O3",
		"-s",
		"-o", c.output,
		c.output+".c")

	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr

	err := gcc.Run()
	return err
}

// Generate takes the specified input-string and writes it as a compiled
// binary to the named output-path.
//
// We generate a temporary file, write our C-source to that and then
// compile via gcc.
func (c *GeneratorC) Generate(input string, output string) error {

	//
	// Save the input and output path away.
	//
	c.input = input
	c.output = output

	//
	// Generate our output program
	//
	err := c.generateSource()
	if err != nil {
		return err
	}

	//
	// Compile it
	//
	err = c.compileSource()
	if err != nil {
		return err
	}

	//
	// Cleanup our source file?  Or leave it alone
	// and output the path of the source-file we generated.
	//
	clean := os.Getenv("CLEANUP")
	if clean == "1" {
		os.Remove(c.output + ".c")
	} else {

		fmt.Printf("generated source file at %s\n", c.output+".c")
	}

	return nil
}

// Register our back-end
func init() {
	Register("c", func() Generator {
		return &GeneratorC{}
	})
}
