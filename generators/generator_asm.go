package generators

// GeneratorASM is a generator that will produce an x86-64 assembly-language
// version of the specified input-program.
//
// The assembly language file will be compiled by nasm.
type GeneratorASM struct {
}

// Generate takes the specified input-string and writes it as a compiled
// binary to the named output-path.
//
// We generate a temporary file, write our assembly language file to that
// and then compile via nasm, and link with ld.
func (g *GeneratorASM) Generate(input string, output string) error {
	return nil
}

// Register our back-end
func init() {
	Register("asm", func() Generator {
		return &GeneratorASM{}
	})
}
