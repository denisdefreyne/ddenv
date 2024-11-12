.PHONY: build
build:
	goreleaser build --snapshot --clean

.PHONY: install
install: build
	install -v dist/ddenv_darwin_arm64_v8.0/ddenv ~/bin
