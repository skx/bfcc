# BrainFuck Compiler Challenge

The aim of this repository is to contain a BrainFuck compiler, written in Golang, in less than a day.

## Rough Plan

There are a lot of different ways to skin this cat, but my starting plan is to do this:

* Parse a valid program.
* Generate C-code which corresponds to that input.
* Compile it with GCC.

The next step is to drop the C-generation, and instead generate assembly language:

* Parse a valid program.
* Generate an x86 assembly version of the input.
* Compile it with nasm

The final step would be to drop the use of the assembler, but I suspect I might not get that far.


