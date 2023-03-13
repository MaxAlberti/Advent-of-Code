build:
	GOARCH=wasm GOOS=js go build -C ui/server -o bin/web/app.wasm
	go build -C ui/server -o bin/ui

run: build
	cd ui/server/bin/ && ./ui