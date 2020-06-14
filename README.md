# BrainFuck Compiler Challenge

The aim of this repository is to contain a BrainFuck compiler, written in Golang, in less than a day.

[Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) is an esoteric programming language created in 1993 by Urban Müller, and is notable for its extreme minimalism.  It supports only a few instructions, and is practically unreadable.

That said brainfuck, despite the name, has a good history in the programming world as being something simple and well-defined to play with.



## Usage

You can install the compiler via:

    $ go get github.com/skx/bfcc

Once installed execute the compiler as follows to produce the default executable at `./a.out`:

    $ bfcc ./examples/mandelbrot.bf

Rather than compile, then run, you can add `-run` to your invocation:

    $ bfcc -run ./examples/bizzfuzz.bf

Finally if you prefer you can specify an output name for the compiled result:

    $ bfcc [-run] ./examples/bizzfuzz.bf ./bf


**NOTE**: You need to have


## Rough Plan

There are a lot of different ways to skin this cat, but my starting plan was to do this:

* [x] Parse a valid program.
* [x] Generate C-code which corresponds to that input.
* [x] Compile it with `gcc`.

Once that was completed the next step would be to drop the use of GCC and instead generate assembly language:

* [x] Parse a valid program.
* [x] Generate an x86 assembly version of the input.
* [x] Compile it with `nasm`, and link with `ldd`.

Annoyingly I noticed that when I compiled the Mandelbrot example via the assembly language version it was slower than the C-based version.  So I had some more things to do:

* [x] Optimize the generated assembly language.



## Timeline

* Project occurred to me overnight.
* Started work at 12:00.
  * Generated initial version of BrainFuck -> C in 30 minutes.
* Paused for a break at 12:58
  * Had added documentation, added more sample programs, and added test-suite.
* Started work at 15:00 again.
  * Implemented trivial assembly language version by 15:30.

You can walk backwards in the commit-history if you wish, but the final version of the C-generating version was:

* [cadb19d6c75a5febde56f53423a9668ee8f6bd25](https://github.com/skx/bfcc/tree/cadb19d6c75a5febde56f53423a9668ee8f6bd25)


## Test Programs

There are a small collection of test-programs located beneath the [examples/](examples/) directory.

Each example has a `.bf` suffix, and there is a corresponding output file for each input to show the expected output.

You can run `make test` to run all the scripts, and compare their generated output with the expected result.


## Speed

I've implemented two simple "compilers":

* The first generated C code.
  * This was then compiled via `gcc` (with `-O3`).
* The second generated assembly language code.
  * Compile via `nasm`, and linked with `ld`.

The most complicated program I've run was the Mandelbrot generator, and surprisingly the runtime of the C-based version is faster:

| Version  | RunTime |
|----------|---------|
| C        | 1.177s  |
| Assembly | 2.694s  |


Of course there are obvious optimizations to be made, which is why I structured the assembly language output as I did.  For example `>` is used to increase our index.  I compile `>` to `add r8, 1` as the R8 register is used for our index.

However I compile "`>>>`" into "`add r8, 1; add r8, 1; add r8, 1` when instead I should compile it into `add r8,3`.  (i.e. I can collapse multiple identical instances of the increase/decrease instructions into a single instruction.)

The starting version of the assembly output was :

* [aebb14ccb548a2249bc32bb1f82fe9070518cc3c](https://github.com/skx/bfcc/tree/aebb14ccb548a2249bc32bb1f82fe9070518cc3c)

The optimized version was :

* [91d6712bcb4b41e9fd963f60da2753d62ee789d1](https://github.com/skx/bfcc/commit/91d6712bcb4b41e9fd963f60da2753d62ee789d1)

With that the timings become:

| Version            | RunTime |
|--------------------|---------|
| C                  | 1.177s  |
| Assembly           | 2.694s  |
| Optimized Assembly | 1.542s  |


## Bug Reports?

Please [file an issue](https://github.com/skx/bfcc/issues)


Steve
--
