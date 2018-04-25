MAIN_VERSION:=$(shell git describe --abbrev=0 --tags || echo "0.1")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v -e /vendor/ -e github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces)
LDFLAGS:=-ldflags "-X github.com/ch-robinson/vault-elastic-plugin/plugin.Version=${VERSION}"

default: test
	make build 
	
test:
	@echo "mode: count" > coverage-all.out
	@echo
	@echo "************** SKIPPING TESTS IN PACKAGES ****************"
	@echo
	@echo "skipping github.com/ch-robinson/vault-elastic-plugin/plugin/interfaces"
	@echo
	@echo "**********************************************************"
	@echo

	@$(foreach pkg,$(PACKAGES), \
		go test -p=1 -cover -covermode=count -coverprofile=coverage.out ${pkg} || exit 1; \
		tail -n +2 coverage.out >> coverage-all.out;)

	@COVERAGE=$$(go tool cover -func=coverage-all.out | tail -1 | tr -d '[:space:]' | tr -d '()' | tr -d '%' | tr -d ':' | sed -e 's/total//g' | sed -e 's/statements//g'); \
	echo "Total Coverage: $${COVERAGE}"; 

.PHONY: test

cover: test
	@go tool cover -html=coverage-all.out

.PHONY: cover

depends:
	glide up

.PHONY: depends
	
run:
	go run main.go

.PHONY: run

build:
	GOOS=linux GOARCH=amd64 go build -a -o bin/linux/vault-elastic-plugin-x86-64 main.go
	GOOS=windows GOARCH=amd64 go build -a -o bin/windows/vault-elastic-plugin-x86-64.exe main.go
	GOOS=darwin GOARCH=amd64 go build -a -o bin/mac/vault-elastic-plugin-x86-64 main.go

.PHONY: build 

clean:
	rm -rf vendor bin coverage.out coverage-all.out

.PHONY: clean
