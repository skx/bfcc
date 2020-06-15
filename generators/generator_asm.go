package generators

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/skx/bfcc/lexer"
)

// GeneratorASM is a generator that will produce an x86-64 assembly-language
// version of the specified input-program.
//
// The assembly language file will be compiled by nasm.
type GeneratorASM struct {

	// input source
	input string

	// file to write to
	output string
}

// generateSource produces a version of the program as X86-64 assembly language.
func (g *GeneratorASM) generateSource() error {
	var buff bytes.Buffer
	var programStart = `
global _start
section .text

_start:
  mov r8, stack
`
	buff.WriteString(programStart)

	//
	// Keep track of "[" here.
	//
	// These are loop opens.
	//
	opens := []int{}

	//
	// Create a lexer for the input program
	//
	l := lexer.NewLexer(g.input)

	//
	// Loop forever, processing the next token
	//
	tok := l.Next()

	//
	// We keep track of the loop-labels here.
	//
	// Each time we see a new loop-open "[" we bump this
	// by one.
	//
	i := 0

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
			buff.WriteString(fmt.Sprintf("  add r8, %d\n", tok.Repeat))

		case lexer.LESS:
			buff.WriteString(fmt.Sprintf("  sub r8, %d\n", tok.Repeat))

		case lexer.PLUS:
			buff.WriteString(fmt.Sprintf("  add byte [r8], %d\n", tok.Repeat))

		case lexer.MINUS:
			buff.WriteString(fmt.Sprintf("  sub byte [r8], %d\n", tok.Repeat))

		case lexer.OUTPUT:
			buff.WriteString("  mov rax, 1\n")  // SYS_WRITE
			buff.WriteString("  mov rdi, 1\n")  // STDOUT
			buff.WriteString("  mov rsi, r8\n") // data-comes-here
			buff.WriteString("  mov rdx, 1\n")  // one byte
			buff.WriteString("  syscall\n")     // Syscall

		case lexer.INPUT:
			buff.WriteString("  mov rax, 0\n")  // SYS_READ
			buff.WriteString("  mov rdi, 0\n")  // STDIN
			buff.WriteString("  mov rsi, r8\n") // Dest
			buff.WriteString("  mov rdx, 1\n")  // one byte
			buff.WriteString("  syscall\n")     // syscall

		case lexer.LOOP_OPEN:

			//
			// Open of a block.
			//
			// If the index-value is zero then jump to the
			// end of the while-loop.
			//
			// NOTE: We repeat the test at the end of the
			// loop so the label here is AFTER our condition
			//
			i++
			buff.WriteString("  cmp byte [r8], 0\n")
			buff.WriteString(fmt.Sprintf("  je close_loop_%d\n", i))
			buff.WriteString(fmt.Sprintf("label_loop_%d:\n", i))
			opens = append(opens, i)

		case lexer.LOOP_CLOSE:

			// "]" can only follow an "[".
			//
			// Every time we see a "[" we save the ID onto a
			// temporary stack.  So we're gonna go back to the
			// most recent open.
			//
			// This will cope with nesting.
			//
			if len(opens) < 1 {
				fmt.Printf("close before open.  bug?  bogus program?\n")
				os.Exit(1)
			}

			//
			// Get the last label-ID
			//
			last := opens[len(opens)-1]

			//
			// Remove it from our list now.
			//
			opens = opens[:len(opens)-1]

			//
			// What we could do here is jump back to the
			// start of our loop.
			//
			// The test would be made, and if it failed we'd
			// end up back at the end of the loop.
			//
			// However we're tricksy hobbitses, so we run
			// the test again, and only jump back if the
			// loop is not yet over.
			//
			// As per suggestion from Wikipedia.
			//
			// This has a cost of comparing twice, but
			// a benefit of ensuring we don't jump more than
			// we need to.
			//
			// NOTE: That we jump AFTER the conditional
			// test at the start of the loop, because
			// running it twice would be pointless.
			//
			buff.WriteString("  cmp byte [r8], 0\n")

			buff.WriteString(fmt.Sprintf("  jne label_loop_%d\n", last))
			buff.WriteString(fmt.Sprintf("close_loop_%d:\n", last))

		default:
			fmt.Printf("token not handled: %v\n", tok)
			os.Exit(1)
		}

		//
		// Keep processing
		//
		tok = l.Next()
	}

	// terminate
	buff.WriteString("  mov rax, 60\n")
	buff.WriteString("  mov rdi, 0\n")
	buff.WriteString("  syscall\n")

	// program-area
	buff.WriteString("section .bss\n")
	buff.WriteString("stack: resb 300000\n")

	// Output to a file
	err := ioutil.WriteFile(g.output+".s", buff.Bytes(), 0644)
	return err
}

func (g *GeneratorASM) compileSource() error {

	// nasm to compile to object-code
	nasm := exec.Command("nasm", "-f", "elf64", "-o", g.output+".o", g.output+".s")
	nasm.Stdout = os.Stdout
	nasm.Stderr = os.Stderr

	err := nasm.Run()
	if err != nil {
		return err
	}

	// ld to link to an executable
	ld := exec.Command("ld", "-m", "elf_x86_64", "-o", g.output, g.output+".o")
	ld.Stdout = os.Stdout
	ld.Stderr = os.Stderr

	err = ld.Run()
	if err != nil {
		return err
	}

	return nil
}

// Generate takes the specified input-string and writes it as a compiled
// binary to the named output-path.
//
// We generate a temporary file, write our assembly language file to that
// and then compile via nasm, and link with ld.
func (g *GeneratorASM) Generate(input string, output string) error {

	//
	// Save the input and output path away.
	//
	g.input = input
	g.output = output

	//
	// Generate our assembly
	//
	err := g.generateSource()
	if err != nil {
		return err
	}

	//
	// Compile it
	//
	err = g.compileSource()
	if err != nil {
		return err
	}

	clean := os.Getenv("CLEANUP")
	if clean == "1" {
		os.Remove(g.output + ".s")
		os.Remove(g.output + ".o")
	} else {

		os.Remove(g.output + ".o")
		fmt.Printf("generated source file at %s\n", g.output+".s")
	}

	return nil
}

// Register our back-end
func init() {
	Register("asm", func() Generator {
		return &GeneratorASM{}
	})
}
