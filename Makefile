.PHONY: all clean prebuild

VERSION := latest

default: bin/netsoc

prebuild:
	cat tools.go | sed -nr 's|^\t_ "(.+)"$$|\1|p' | xargs -tI % go get %
	protoc --experimental_allow_proto3_optional --go_out=. converter/object.proto

bin/%: prebuild
	go build -o $@ ./cmd/$(shell basename $@)

dev/%: prebuild
	$(eval BIN = $(shell basename $@))
	CompileDaemon -exclude-dir=.git -build="go build -o bin/$(BIN) ./cmd/$(BIN)" \
		-include '*.vs' -include '*.fs'

clean:
	-rm -f bin/*
