# BrainFuck Compiler Challenge

The aim of this repository is to contain a BrainFuck compiler, written in Golang, in less than a day.

[Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) is an esoteric programming language created in 1993 by Urban Müller, and is notable for its extreme minimalism.  It supports only a few instructions, and is practically unreadable.

That said brainfuck, despite the name, as a good history in the programming world as being something simple and well-defined to play with.



## Rough Plan

There are a lot of different ways to skin this cat, but my starting plan is to do this:

* [x] Parse a valid program.
* [x] Generate C-code which corresponds to that input.
* [x] Compile it with `gcc`.

Now that this is done the next step would be to drop the use of GCC and instead generate assembly language:

* [x] Parse a valid program.
* [x] Generate an x86 assembly version of the input.
* [x] Compile it with `nasm`, and link with `ldd`.

A completely-final step would be to drop the use of the assembler entirely, generating a native ELF binary, but I suspect I might not get that far.



## Timeline

* Project occurred to me overnight.
* Started work at 12:00.
  * Generated initial version of BrainFuck -> C in 30 minutes.
* Paused for a break at 12:58
  * Had added documentation, added more sample programs, and added test-suite.
* Started work at 15:00 again.
  * Implemented trivial assembly language version.


## Test Programs

There are a small collection of test-programs located beneath the [examples/](examples/) directory.

Each example has a `.bf` suffix, and there is a corresponding output file for each input to show the expected output.

You can run `make test` to run all the scripts, and compare their generated output with the expected result.


## Usage

You can install the compiler via:

    $ go get github.com/skx/bfcc

Once installed execute the compiler like so:

    $ bfcc ./examples/mandelbrot.bf
    $ time ./a.out
    ..
    real	0m1.177s
    user	0m1.172s
    sys	0m0.000s
    $

Rather than compile, then run, you can add `-run` to your invocation:

    $ bfcc -run ./examples/bizzfuzz.bf

Finally if you prefer you can specify an output name for the compiled result:

    $ bfcc [-run] ./examples/bizzfuzz.bf ./bf


## Bug Reports?

Please [file an issue](https://github.com/skx/bfcc/issues)


Steve
--
