GOPATH:=$(shell go env GOPATH)
VERSION:=$(shell git describe --tags --always)
#PATH:=$(shell $(PATH)):$(shell $(GOPATH))/bin

ifeq ($(go env GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_CONFIG_FILES=$(shell $(Git_Bash) -c "find internal/conf -name *.proto")
else ifeq ($(go env GOHOSTOS), darwin)
	PATH:=$(PATH):$(GOPATH)/bin
	SHELL=env PATH=$(PATH) /bin/bash
	INTERNAL_CONFIG_FILES=$(shell find internal/conf -name *.proto)
else
	PATH:=$(PATH):$(GOPATH)/bin
	SHELL=env PATH=$(PATH) /bin/bash
	INTERNAL_CONFIG_FILES=$(shell find internal/conf -name *.proto)
endif


.PHONY: env
# initilize env
env:
	export GOPROXY=https://goproxy.cn
	export GOSUMDB="off"


.PHONY: init
# init env
init:
	go get github.com/google/wire/cmd/wire@v0.5.0
	go install github.com/codeskyblue/fswatch@latest


.PHONY: config
# generate internal proto
config:
	protoc --proto_path=. \
 	       --go_out=paths=source_relative:. \
	       $(INTERNAL_CONFIG_FILES)


.PHONY: initConfig
# initilize a config file
initConfig:
	mkdir -p ./configs && cp internal/conf/config.dev.yaml ./configs/config.yaml


.PHONY: generate
# generate config & wire_gen
generate:
	./wireGenerate.sh
	go generate ./cmd/...


.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/...

.PHONY: buildApi
# buildApi
buildApi:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/api

.PHONY: buildConsumer
# buildConsumer
buildConsumer:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/consumer

.PHONY: buildScript
# buildScript
buildScript:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./cmd/script


.PHONY: runApi
# start api server
runApi:
	fswatch --config cmd/api/.fsw.yml

.PHONY: runConsumer
# start consumer
runConsumer:
	fswatch --config cmd/consumer/.fsw.yml

.PHONY: runScript
# start script
runScript:
	fswatch --config cmd/script/.fsw.yml


.PHONY: all
# generate all
all:
	make env;
	make init;
	make config;
	make initConfig;
	make generate;
	cd proto && make all;
	make build;
	make runApi


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
