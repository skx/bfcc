//
// Trivial "compiler" for BrainFuck
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Given an input brainfuck program generate a comparable assembly version.
//
// Return the filename we generated, or error.
func generateProgram(source string) (string, error) {
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
	l := NewLexer(source)

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
	for tok.Type != EOF {

		//
		// Output different things depending on the token-type
		//
		switch tok.Type {

		case GREATER:
			buff.WriteString(fmt.Sprintf("  add r8, %d\n", tok.Repeat))

		case LESS:
			buff.WriteString(fmt.Sprintf("  sub r8, %d\n", tok.Repeat))

		case PLUS:
			buff.WriteString(fmt.Sprintf("  add byte [r8], %d\n", tok.Repeat))

		case MINUS:
			buff.WriteString(fmt.Sprintf("  sub byte [r8], %d\n", tok.Repeat))

		case OUTPUT:
			buff.WriteString("  mov rax, 1\n")  // SYS_WRITE
			buff.WriteString("  mov rdi, 1\n")  // STDOUT
			buff.WriteString("  mov rsi, r8\n") // data-comes-here
			buff.WriteString("  mov rdx, 1\n")  // one byte
			buff.WriteString("  syscall\n")     // Syscall

		case INPUT:
			buff.WriteString("  mov rax, 0\n")  // SYS_READ
			buff.WriteString("  mov rdi, 0\n")  // STDIN
			buff.WriteString("  mov rsi, r8\n") // Dest
			buff.WriteString("  mov rdx, 1\n")  // one byte
			buff.WriteString("  syscall\n")     // syscall

		case LOOP_OPEN:

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

		case LOOP_CLOSE:

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

	tmpfile, err := ioutil.TempFile("", "bfcc*.s")
	if err != nil {
		return "", err
	}

	if _, err := tmpfile.Write(buff.Bytes()); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}

	return tmpfile.Name(), nil
}

func main() {

	//
	// Parse command-line flags
	//
	compile := flag.Bool("compile", true, "Compile the assembly file, after generation.")
	cleanup := flag.Bool("cleanup", true, "Remove the temporary assembly file after creation.")
	run := flag.Bool("run", false, "Run the program after compiling.")
	flag.Parse()

	//
	// Get the input filename
	//
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: bfcc [flags] input.file.bf [outfile]\n")
		return
	}

	//
	// Input and output files
	//
	input := flag.Args()[0]
	output := "a.out"
	if len(flag.Args()) == 2 {
		output = flag.Args()[1]
	}

	//
	// Read the input program.
	//
	prog, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("failed to read input file %s: %s\n", input, err.Error())
		return
	}

	//
	// "Compile"
	//
	path, err := generateProgram(string(prog))
	if err != nil {
		fmt.Printf("error writing output: %s\n", err.Error())
		return
	}

	//
	// Compile
	//
	if *compile {

		// nasm to compile to object-code
		nasm := exec.Command("nasm", "-f", "elf64", "-o", fmt.Sprintf("%s.o", output), path)
		nasm.Stdout = os.Stdout
		nasm.Stderr = os.Stderr

		err = nasm.Run()
		if err != nil {
			fmt.Printf("Error launching nasm: %s\n", err)
			return
		}

		// ld to link to an executable
		ld := exec.Command("ld", "-m", "elf_x86_64", "-o", output, fmt.Sprintf("%s.o", output))
		ld.Stdout = os.Stdout
		ld.Stderr = os.Stderr

		err = ld.Run()
		if err != nil {
			fmt.Printf("Error launching ld: %s\n", err)
			return
		}

	}

	if *run {
		exe := exec.Command(output)
		exe.Stdin = os.Stdin
		exe.Stdout = os.Stdout
		exe.Stderr = os.Stderr
		err = exe.Run()
		if err != nil {
			fmt.Printf("Error launching %s: %s\n", output, err)
			os.Exit(1)
		}

	}

	//
	// Cleanup
	//
	if *cleanup {
		os.Remove(path)
	} else {
		fmt.Printf("Generated output left at %s\n", path)
	}
}
