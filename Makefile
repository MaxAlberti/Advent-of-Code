build:
	go build -C cmd/AoC -o ../../bin/AoC

run: build
	./bin/AoC