build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	go build -o bin/ui

run: build
	./bin/ui