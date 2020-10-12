.PHONY: all clean phony_explicit

VERSION := latest

default: bin/netsoc

phony_explicit:

bin/%: phony_explicit
	go build -o $@ ./cmd/$(shell basename $@)

dev/%: phony_explicit
	cat tools.go | sed -nr 's|^\t_ "(.+)"$$|\1|p' | xargs -tI % go get %

	$(eval BIN = $(shell basename $@))
	CompileDaemon -exclude-dir=.git -build="go build -o bin/$(BIN) ./cmd/$(BIN)" \
		-include '*.vs' -include '*.fs'

clean:
	-rm -f bin/*
