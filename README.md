# BrainFuck Compiler Challenge

The aim of this repository is to contain a BrainFuck compiler, written in Golang, in less than a day.

## Rough Plan

There are a lot of different ways to skin this cat, but my starting plan is to do this:

* [x] Parse a valid program.
* [x] Generate C-code which corresponds to that input.
* [x] Compile it with GCC.

Now that this is done the next step would be to drop the use of GCC and instead generate assembly language:

* Parse a valid program.
* Generate an x86 assembly version of the input.
* Compile it with nasm/gcc/similar.

A completely-final step would be to drop the use of the assembler entirely, generating a native ELF binary, but I suspect I might not get that far.

## Timeline

* Plan occurred to me overnight
* Started work at 12:00.
* Generated initial version of BrainFuck -> C in 30 minutes.

Time taken: 30 minutes.
