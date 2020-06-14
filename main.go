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

// Given an input brainfuck program generate a comparable C version.
//
// Return the filename we generated, or error
func generateProgram(source string) (string, error) {
	var buff bytes.Buffer
	var programStart = `
	extern int putchar(int);
	extern char getchar();

	char array[30000];
	int idx = 0;

	int main (int arc, char *argv[]) {
	`
	buff.WriteString(programStart)

	bts := []byte(source)
	for _, bt := range bts {
		switch bt {
		case '>':
			buff.WriteString("idx++;\n")
		case '<':
			buff.WriteString("idx--;\n")
		case '+':
			buff.WriteString("array[idx]++;\n")
		case '-':
			buff.WriteString("array[idx]--;\n")
		case '.':
			buff.WriteString("putchar(array[idx]);\n")
		case ',':
			buff.WriteString("array[idx] = getchar();\n")
		case '[':
			buff.WriteString("while (array[idx]) {\n")
		case ']':
			buff.WriteString("}\n")
		}
	}

	buff.WriteString("}\n")

	tmpfile, err := ioutil.TempFile("", "bfcc*.c")
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
	compile := flag.Bool("compile", true, "Compile the generated C, after generation")
	run := flag.Bool("run", false, "Run the program after compiling")
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
		output = os.Args[1]
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
		gcc := exec.Command("gcc", "-O3", "-Ofast", "-o", output, path)
		gcc.Stdout = os.Stdout
		gcc.Stderr = os.Stderr

		err = gcc.Run()
		if err != nil {
			fmt.Printf("Error launching gcc: %s\n", err)
			return
		}
	}

	if *run {
		exe := exec.Command(output)
		exe.Stdout = os.Stdout
		exe.Stderr = os.Stderr
		err = exe.Run()
		if err != nil {
			fmt.Printf("Error launching %s: %s\n", output, err)
			os.Exit(1)
		}

	}
}
