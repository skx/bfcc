[![Go Report Card](https://goreportcard.com/badge/github.com/skx/bfcc)](https://goreportcard.com/report/github.com/skx/bfcc)
[![license](https://img.shields.io/github/license/skx/bfcc.svg)](https://github.com/skx/bfcc/blob/master/LICENSE)

Table of Contents
=================

* [BrainFuck Compiler Challenge](#brainfuck-compiler-challenge)
   * [Usage](#usage)
   * [My Approach](#my-approach)
   * [Timeline](#timeline)
      * [Debugging the generated program](#debugging-the-generated-program)
   * [Test Programs](#test-programs)
   * [Future Plans?](#future-plans)
   * [Bug Reports?](#bug-reports)



# BrainFuck Compiler Challenge

The aim of this repository was to contain a BrainFuck compiler, written in Golang, and completed in less than a day.

[Brainfuck](https://en.wikipedia.org/wiki/Brainfuck) is an esoteric programming language created in 1993 by Urban Müller, and is notable for its extreme minimalism.  It supports only a few instructions, and is practically unreadable.

That said brainfuck, despite the name, has a good history in the programming world as being something simple and well-defined to play with.



## Usage

You can install the compiler via:

    $ go get github.com/skx/bfcc

Once installed execute the compiler as follows to produce the default executable at `./a.out`:

    $ bfcc ./examples/mandelbrot.bf
    $ ./a.out

Rather than compile, then run, you can add `-run` to your invocation:

    $ bfcc -run ./examples/bizzfuzz.bf

Finally if you prefer you can specify an output name for the compiled result:

    $ bfcc [-run] ./examples/bizzfuzz.bf ./bf
    $ ./bf

There are two backends included, one generates an assembly language source-file, and compiles with with `nasm`, and the other generates C-code, which is compiled via `gcc`.

By default the assembly-language backend is selected, because this is the thing that I was more interested in writing, due to that you **must** have `nasm` installed to compile the generated assembly language file.

Use the `-backend` flag to specify the backend which you prefer to use:

    $ bfcc -backend=c   ./examples/mandelbrot.bf ./mb-c
    $ bfcc -backend=asm ./examples/mandelbrot.bf ./mb-asm

You'll see slightly difference sizes in the two executable:

    $ ls -lash mb-c mb-asm
    76K -rwxr-xr-x 1 skx skx 73K Jun 15 10:11 mb-asm
    36K -rwxr-xr-x 1 skx skx 34K Jun 15 10:11 mb-c

But both should work identically; if they do not that's a bug in the generated C/assembly source files I've generated!




## My Approach

In the end it took me about four hours to get something I was happy with, and later I've improved it a little more:

* Initially I generated C-code, [as you can see here](https://github.com/skx/bfcc/blob/cadb19d6c75a5febde56f53423a9668ee8f6bd25/main.go).
  * This code was largely inspired by the [wikipedia brainfuck page](https://en.wikipedia.org/wiki/Brainfuck)
  * The C-code was compiled by GCC to produce an executable.
* Then I started generating assembly-language code, which looked [something like this](https://github.com/skx/bfcc/blob/aebb14ccb548a2249bc32bb1f82fe9070518cc3c/main.go).
  * The generated assembly was compiled by `nasm`, and linked with `ld` to produce an executable.
* Once the assembly language code was working I optimized it.
  * Collapsing multiple identical instructions to one.
  * Improving the way loop-handling was generatoed.
* Finally I cleaned up and improved the code.
  * Implementing a separate lexer.
  * Allowing the use of pluggable backends, so we could generate both C and Assembly Language output (but only one at a time).




## Timeline

* Project occurred to me overnight.
* Started work at 12:00.
  * Generated initial version of BrainFuck -> C in 30 minutes.
* Paused for a break at 12:58
  * Had added documentation, added more sample programs, and added test-suite.
* Started work at 15:00 again.
  * Implemented trivial assembly language version by 15:30.
* Spent another hour cleaning up comments, _this_ README.md file, and applying basic optimizations.
* After that I slowly made improvements
  * Adding a lexer in [#4](https://github.com/skx/bfcc/pull/4)
  * Allowing the generation of either C or assembly in [#6](https://github.com/skx/bfcc/pull/6)
  * Allow generating a breakpoint instruction when using the assembly-backend in [#7](https://github.com/skx/bfcc/pull/7).


### Debugging the generated program

If you run the compiler with the `-debug` flag a breakpoint will be generated
immediately at the start of the program.  You can use that breakpoint to easily
debug the generated binary via `gdb`.

    $ bfcc -debug ./examples/hello-world.bf

Now you can launch that binary under `gdb`, and run it:

    $ gdb ./a.out
    (gdb) run
    ..
    Program received signal SIGTRAP, Trace/breakpoint trap.
    0x00000000004000bb in _start ()

Disassemble the code via `disassemble`, and step over instructions one at a time via `stepi`.  If your program is long you might see a lot of output from the `disassemble` step.

    (gdb) disassemble
    Dump of assembler code for function _start:
       0x00000000004000b0 <+0>:	movabs $0x600290,%r8
       0x00000000004000ba <+10>:	int3
    => 0x00000000004000bb <+11>:	addb   $0x8,(%r8)
       0x00000000004000bf <+15>:	cmpb   $0x0,(%r8)
       0x00000000004000c3 <+19>:	je     0x40013f <close_loop_1>
    End of assembler dump.

You can set a breakpoint at a line in the future, and continue running till
you hit it, with something like this:

     (gdb) break *0x00000000004000c3
     (gdb) cont

Once there inspect the registers with commands like these two:

     (gdb) print $r8
     (gdb) print *$r8
     (gdb) info registers

Further documentation can be found in the `gdb` manual, which is worth reading
if you've an interest in compilers, debuggers, and decompilers.



## Test Programs

There are a small collection of test-programs located beneath the [examples/](examples/) directory.

Each example has a `.bf` suffix, and there is a corresponding output file for each input to show the expected output.

You can run `make test` to run all the scripts, and compare their generated output with the expected result.




## Future Plans?

Mostly none.

It might be cute to convert the assembly, such that `gcc` could compile it.  That would drop the `nasm` dependency, but it's not a big deal.  Patches welcome if you want to have a stab at it.

Otherwise more backends might be nice, but I guess the two existing ones are the most obvious.  Due to the way the code is structured adding a new one would be trivial though.




## Bug Reports?

Please [file an issue](https://github.com/skx/bfcc/issues)


Steve
--
