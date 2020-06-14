//
// Trivial "compiler" for BrainFuck
//
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

// Given an input brainfuck program generate a comparable C version.
func compile(source string) string {
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
	return buff.String()
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("Usage: bfcc <infile> [<outfile>]\n")
		return
	}

	input := os.Args[1]
	output := "out.c"
	if len(os.Args) == 3 {
		output = os.Args[2]
	}

	// Read input
	prog, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("Failed to read input file %s: %s\n", input, err.Error())
		return
	}

	// "Compile"
	out := compile(string(prog))

	// Write output
	err = ioutil.WriteFile(output, []byte(out), 0644)
	if err != nil {
		fmt.Printf("failed to write to %s: %s\n", output, err.Error())
		return
	}

}
