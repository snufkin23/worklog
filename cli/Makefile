APP := wl

.PHONY: run build test fmt vet tidy

run:
	go run ./cmd/wl

build:
	go build -o bin/$(APP) ./cmd/wl

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

tidy:
	go mod tidy
