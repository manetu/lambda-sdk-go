
all: build

main.wasm: main.go Makefile
	tinygo build -o $@ --target=wasi $<

build: main.wasm
