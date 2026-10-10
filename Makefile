BIN = build/emulator

build:
	go build -o $(BIN) ./src

run: build
	./$(BIN) $(ARGS)

test:
	go test -v ./tests/...

clean:
	rm -rf build

.PHONY: build run test clean
