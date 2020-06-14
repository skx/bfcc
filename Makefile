
build:
	go build .

test: build
	@echo "hello-world test"
	./bfcc -run examples/hello-world.bf > x
	diff examples/hello-world.out x && rm x || echo "Failed Hello-World"

	@echo "fibonacci test"
	./bfcc -run examples/fibonacci.bf > x
	diff examples/fibonacci.out x && rm x || echo "Failed Fibonacci"

	@echo "bizzfuzz test"
	./bfcc -run examples/bizzfuzz.bf > x
	diff examples/bizzfuzz.out x && rm x || echo "Failed BizzFuzz"

	@echo "quine test"
	./bfcc -run examples/quine.bf > x
	diff examples/quine.out x && rm x || echo "Failed Quine"

	@echo "mandelbrot test"
	./bfcc -run examples/mandelbrot.bf > x
	diff examples/mandelbrot.out x && rm x || echo "Failed Mandelbrot"
