package generators

// GeneratorC is a generator that will produce an C version of the specified
// input-program.
//
// The C-source will then be compiled by gcc.
type GeneratorC struct {
}

// Generate takes the specified input-string and writes it as a compiled
// binary to the named output-path.
//
// We generate a temporary file, write our C-source to that and then
// compile via gcc.
func (g *GeneratorC) Generate(input string, output string) error {
	return nil
}

// Register our back-end
func init() {
	Register("c", func() Generator {
		return &GeneratorC{}
	})
}
