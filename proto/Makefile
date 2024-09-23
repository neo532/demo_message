GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
PATH:=$(PATH):$(GOPATH)/bin
SHELL:=env PATH=$(PATH) $(SHELL)

ifeq ($(GOHOSTOS), windows)
	#the `find.exe` is different from `find` in bash/shell.
	#to see https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/find.
	#changed to use git-bash.exe to run find cli or other cli friendly, caused of every developer has a Git.
	#Git_Bash= $(subst cmd\,bin\bash.exe,$(dir $(shell where git)))
	Git_Bash=$(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	API_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.api.proto")
	CLIENT_PROTO_FILES=$(shell $(Git_Bash) -c "find client -name *.api.proto")
	ERR_PROTO_FILES=$(shell $(Git_Bash) -c "find api -name *.err.proto")
else
	API_PROTO_FILES=$(shell find api -name *.api.proto)
	CLIENT_PROTO_FILES=$(shell find client -name *.api.proto)
	ERR_PROTO_FILES=$(shell find api -name *.err.proto)
endif

.PHONY: init
# init env
init:
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install github.com/favadi/protoc-go-inject-tag@latest
	go install github.com/neo532/apitool@latest


.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
		   --validate_out=lang=go,paths=source_relative:./api \
	       $(API_PROTO_FILES)


.PHONY: client
# generate client proto
client:
	apitool httpclient $(CLIENT_PROTO_FILES)


.PHONY: error
# generate err proto
error:
	protoc --proto_path=. \
		--proto_path=./third_party \
		--go_out=paths=source_relative:. \
		--go-errors_out=paths=source_relative:. \
		$(ERR_PROTO_FILES)


.PHONY: all
# generate all
all:
	make api;
	make client;
	make error;

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
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
