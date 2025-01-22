.DEFAULT_GOAL := goapp

.PHONY: all
all: clean test bench goapp

.PHONY: goapp
goapp:
	mkdir -p bin
	go build -o bin ./...

.PHONY: clean
clean:
	go clean
	rm -f bin/*

.PHONY: test
test:
	go test -cover ./...

.PHONY: bench
bench:
	go test -bench=. -benchmem ./...
