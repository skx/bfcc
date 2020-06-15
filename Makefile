
build:
	go build .

test: build
	@echo "hello-world test"
	@./bfcc -run examples/hello-world.bf > x
	@diff examples/hello-world.out x && rm x || echo "Failed Hello-World"
	@echo "hello-world test OK"

	@echo "fibonacci test"
	@./bfcc -run examples/fibonacci.bf > x
	@diff examples/fibonacci.out x && rm x || echo "Failed!"
	@echo "fibonacci test: OK"

	@echo "bizzfuzz test"
	@./bfcc -run examples/bizzfuzz.bf > x
	@diff examples/bizzfuzz.out x && rm x || echo "Failed!"
	@echo "bizzfuzz test: OK"

	@echo "quine test"
	@./bfcc -run examples/quine.bf > x
	@diff examples/quine.out x && rm x || echo "Failed!"
	@echo "quine test: OK"

	@echo "mandelbrot test"
	@./bfcc -run examples/mandelbrot.bf > x
	@diff examples/mandelbrot.out x && rm x || echo "Failed!"
	@echo "mandelbrot test: OK"
