//
// Trivial "compiler" for BrainFuck
//
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/skx/bfcc/generators"
)

func main() {

	//
	// Parse command-line flags
	//
	backend := flag.String("backend", "asm", "The backend to use for compilation.")
	cleanup := flag.Bool("cleanup", true, "Remove the generated files after creation.")
	run := flag.Bool("run", false, "Run the program after compiling.")
	flag.Parse()

	//
	// Ensure the backend we have is available
	//
	helper := generators.GetGenerator(*backend)
	if helper == nil {

		fmt.Printf("Unknown backend %s - valid backends are:\n", *backend)
		all := generators.Available()
		for _, name := range all {
			fmt.Printf("\t%s\n", name)
		}
		return
	}

	//
	// Ensure we have an input filename
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
	// Read the input program
	//
	prog, err := ioutil.ReadFile(input)
	if err != nil {
		fmt.Printf("failed to read %s: %s\n", input, err.Error())
	}

	//
	// Will we cleanup ?
	//
	if *cleanup {
		os.Setenv("CLEANUP", "1")
	} else {
		os.Setenv("CLEANUP", "0")
	}

	//
	// Generate the compiled version
	//
	err = helper.Generate(string(prog), output)
	if err != nil {
		fmt.Printf("error generating binary: %s\n", err.Error())
		return
	}

	//
	// Are we running the program?  Then do so.
	//
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
}
