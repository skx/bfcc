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
	// Loop forever, processing the next token
	//
	tok := l.Next()

	//
	// We'll process the complete program until
	// we hit an end of file/input
	//
	for tok.Type != lexer.EOF {

		//
		// Output different things depending on the token-type
		//
		switch tok.Type {

		case lexer.GREATER:
			buff.WriteString(fmt.Sprintf("  idx += %d;\n", tok.Repeat))

		case lexer.LESS:
			buff.WriteString(fmt.Sprintf("  idx -= %d;\n", tok.Repeat))

		case lexer.PLUS:
			buff.WriteString(fmt.Sprintf("  array[idx] += %d;\n", tok.Repeat))

		case lexer.MINUS:
			buff.WriteString(fmt.Sprintf("  array[idx] -= %d;\n", tok.Repeat))

		case lexer.OUTPUT:
			buff.WriteString("  putchar(array[idx]);\n")

		case lexer.INPUT:
			buff.WriteString("  array[idx] = getchar();\n")

		case lexer.LOOPOPEN:
			buff.WriteString("  while (array[idx]) {\n")

		case lexer.LOOPCLOSE:
			buff.WriteString("}\n")

		default:
			fmt.Printf("token not handled: %v\n", tok)
			os.Exit(1)
		}

		//
		// Keep processing
		//
		tok = l.Next()
	}

	// Close the main-function
	buff.WriteString("}\n")

	// Output to a file
	err := ioutil.WriteFile(c.output+".c", buff.Bytes(), 0644)
	return err
}

// compileSource uses gcc to compile the generated source-code
func (c *GeneratorC) compileSource() error {

	gcc := exec.Command("gcc", "-O3", "-Ofast", "-o", c.output, c.output+".c")
	gcc.Stdout = os.Stdout
	gcc.Stderr = os.Stderr

	err := gcc.Run()
	if err != nil {
		return err
	}
	return nil
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
	// and output the name of the program we created.
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
