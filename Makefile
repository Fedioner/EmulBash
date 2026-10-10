BIN = build/emulator

build:
	go build -o $(BIN) ./src

run: build vfs
	./$(BIN) $(ARGS)

test:
	go test -v ./tests/...

vfs:
	./scripts/make_vfs.sh

clean:
	rm -rf build

.PHONY: build run test vfs clean
