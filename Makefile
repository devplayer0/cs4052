.PHONY: all clean bin/%

VERSION := latest

default: bin/netsoc

bin/%:
	CGO_ENABLED=0 go build -o $@ ./cmd/$(shell basename $@)

dev/%:
	cat tools.go | sed -nr 's|^\t_ "(.+)"$$|\1|p' | xargs -tI % go get %

	$(eval BIN = $(shell basename $@))
	CompileDaemon -exclude-dir=.git -build="go build -o bin/$(BIN) ./cmd/$(BIN)" \
		-command="bin/$(BIN)" -graceful-kill

clean:
	-rm -f bin/*
