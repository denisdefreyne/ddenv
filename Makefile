.PHONY: build
build:
	go build .

.PHONY: install
install: build
	install -v ddenv ~/bin
