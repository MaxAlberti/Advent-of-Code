SOURCES := $(wildcard pkg/aocyear/*/*)

.PHONY: all

all: build

build_plugins: $(SOURCES)
	@for dir in $^; do \
		echo "Building plugin in $$dir"; \
		go build -buildmode=plugin -o ./$$dir/plugin.so ./$$dir; \
	done

build: build_plugins
	go build -C cmd/AoC -o ../../bin/AoC

run: build
	./bin/AoC