.PHONY: all clean

all: 
	go build -o bin/ ./...

clean:
	@rm -fv bin/*

