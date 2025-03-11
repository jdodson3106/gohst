build: 
	go build -o ./cmd/bin/gohst ./cmd/main.go

run: build
	@./cmd/bin/gohst
