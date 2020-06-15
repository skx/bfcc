// Package generators contains the back-ends to generate an executable
// from an input Brainfuck source.
//
// We ship with two backends by default, one to compile the input program
// to C, and one to compile via assembly language.
package generators

import "sync"

// Generator is the interface which must be implemented by
// a backend to compile our code
type Generator interface {

	// Generate the executable from the given source-file
	Generate(input string, output string) error
}

//
// Everything below here is boilerplate to allow
// class-registration and lookup.
//

// This is a map of known-backends.
var handlers = struct {
	m map[string]NewGenerator
	sync.RWMutex
}{m: make(map[string]NewGenerator)}

// NewGenerator is the signature of a constructor-function to
// instantiate a backend.
type NewGenerator func() Generator

// Register allows a new backend to self-register itself, with a name.
func Register(id string, newfunc NewGenerator) {
	handlers.Lock()
	handlers.m[id] = newfunc
	handlers.Unlock()
}

// GetGenerator will retrieve a generator, by name.
func GetGenerator(id string) (a Generator) {
	handlers.RLock()
	ctor, ok := handlers.m[id]
	handlers.RUnlock()
	if ok {
		a = ctor()
	}
	return
}

// Available returns the names of all registered backend handlers.
func Available() []string {
	var result []string

	// For each handler save the name
	handlers.RLock()
	for index := range handlers.m {
		result = append(result, index)
	}
	handlers.RUnlock()

	// And return the result
	return result
}
