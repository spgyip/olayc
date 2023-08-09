.PHONY: all clean test

all: 
	go build -o bin/ ./...

test:
	go test -v .

clean:
	@rm -fv bin/*

