SHELL := cmd.exe

build:
	@echo "Building the server..."
	go build -o dist\server.exe main.go

run: build
	@echo "Running the server..."
	dist\server.exe

watch:
	@echo "Starting reflex to watch for changes..."
	watchexec -r -e go "make run"
