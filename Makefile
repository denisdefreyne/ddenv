.PHONY: build
build:
	go build ./cmd/ddenv/

.PHONY: install
install: build
	install -v ddenv ~/bin
