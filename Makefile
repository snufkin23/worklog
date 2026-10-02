.PHONY: hooks help server-run server-build cli-build mobile-setup mobile-get mobile-run mobile-analyze test fmt

help:
	@echo "server-run      run the API locally"
	@echo "server-build    build the server binary"
	@echo "cli-build       build the wl binary"
	@echo "mobile-setup    install pinned Flutter via fvm and fetch deps"
	@echo "mobile-get      fetch Flutter dependencies"
	@echo "mobile-run      run the Flutter app"
	@echo "mobile-analyze  run flutter analyze"
	@echo "test            run Go tests for server and cli"
	@echo "fmt             format Go code"

server-run:
	$(MAKE) -C server run

server-build:
	$(MAKE) -C server build

cli-build:
	$(MAKE) -C cli build

mobile-setup:
	$(MAKE) -C mobile setup

mobile-get:
	$(MAKE) -C mobile get

mobile-run:
	$(MAKE) -C mobile run

mobile-analyze:
	$(MAKE) -C mobile analyze

test:
	$(MAKE) -C server test
	$(MAKE) -C cli test

fmt:
	$(MAKE) -C server fmt
	$(MAKE) -C cli fmt

hooks:
	git config core.hooksPath .githooks
	chmod +x .githooks/*

.PHONY: check
check:
	@test -z "$$(gofmt -l server cli)" || (echo "Not formatted, run: make fmt" && gofmt -l server cli && exit 1)
	$(MAKE) -C server vet
	$(MAKE) -C server test
	$(MAKE) -C cli vet
	$(MAKE) -C cli test
	$(MAKE) -C mobile analyze
	$(MAKE) -C mobile test
