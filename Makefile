MAIN_VERSION:=$(shell git describe --always || echo "1.0")
VERSION:=${MAIN_VERSION}\#$(shell git log -n 1 --pretty=format:"%h")
PACKAGES:=$(shell go list ./... | sed -n '1!p' | grep -v -e /vendor/)
LDFLAGS:=-ldflags "-X github.com/ch-robinson/vault-elastic-plugin/main.Version=${VERSION}"

ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
	EXECUTABLE_EXT := .exe
else
    DETECTED_OS := $(shell sh -c 'uname -s 2>/dev/null || echo not')
	EXECUTABLE_EXT :=
endif

LOCAL_VAULT_ROOT_TOKEN := sampleroottoken
PLUGIN_DIRECTORY := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))bin/$(shell echo $(DETECTED_OS) | tr A-Z a-z)


default: test
	make build 
	
test:
	@echo "mode: count" > coverage-all.out
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
	GOOS=darwin GOARCH=amd64 go build -a -o bin/darwin/vault-elastic-plugin-x86-64 main.go

.PHONY: build 

clean:
	rm -rf vendor bin coverage.out coverage-all.out

.PHONY: clean

run-vault:
	@echo "setting up Vault config.hcl ..."
ifeq ($(DETECTED_OS),Windows)
	@echo 'plugin_directory = "$(subst /,\\\\,$(PLUGIN_DIRECTORY))"' > ${PLUGIN_DIRECTORY}/config.hcl
else 
	@echo 'plugin_directory = "$(PLUGIN_DIRECTORY)"' > ${PLUGIN_DIRECTORY}/config.hcl
endif
	vault${EXECUTABLE_EXT} server -dev -dev-root-token-id="${LOCAL_VAULT_ROOT_TOKEN}" -config ${PLUGIN_DIRECTORY}/config.hcl

.PHONY: run-vault

test-plugin: 
ifeq ($(INCLUDE_BUILD),true)
	make -s build 
endif

ifeq ($(ENABLE_VAULT_DB), true)
	@echo "Enabling Vault database"
	@VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=${LOCAL_VAULT_ROOT_TOKEN} vault${EXECUTABLE_EXT} secrets enable database
endif
	@echo "Removing previous plugin"
	@curl --header "X-VAULT-TOKEN:${LOCAL_VAULT_ROOT_TOKEN}" --request DELETE http://127.0.0.1:8200/v1/sys/plugins/catalog/vault-elastic-plugin

	@echo "Registering plugin with Vault"
	@VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=${LOCAL_VAULT_ROOT_TOKEN} vault${EXECUTABLE_EXT} write sys/plugins/catalog/vault-elastic-plugin \
	sha_256=$(shell openssl sha256 $(PLUGIN_DIRECTORY)/vault-elastic-plugin-x86-64$(EXECUTABLE_EXT) | sed 's,SHA256($(PLUGIN_DIRECTORY)/vault-elastic-plugin-x86-64$(EXECUTABLE_EXT))=,,g' | sed -e 's/^[[:space:]]*//') \
	command="vault-elastic-plugin-x86-64${EXECUTABLE_EXT}"

	@echo "Configuring Elastic connection and plugin"
	@VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=${LOCAL_VAULT_ROOT_TOKEN} vault${EXECUTABLE_EXT} write database/config/elastic_test \
	connection_url=${ELASTIC_BASE_URI} \
	username=${ELASTIC_USERNAME} \
	password=${ELASTIC_PASSWORD} \
	plugin_name=vault-elastic-plugin \
	allowed_roles="*"

	@echo "Creating 'my-role'"
	@VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=${LOCAL_VAULT_ROOT_TOKEN} vault${EXECUTABLE_EXT} write database/roles/my-role \
	db_name=elastic_test \
	creation_statements=kibanauser \
	default_ttl=5

	@echo "Running plugin..."
	@# Example success:
	@# {
	@# 	"request_id": "ee9ba65f-465f-a187-0c05-83afe0de1008",
	@# 	"lease_id": "database/creds/my-role/b01dd000-ad88-d617-0480-b9fd7494914e",
	@# 	"lease_duration": 2764800,
	@# 	"renewable": true,
	@# 	"data": {
	@# 		"password": "A1a-7uxq992801vr2wv3",
	@# 		"username": "v-root-my-role-7yxuu1x67wu91q2"
	@# 	},
	@# 	"warnings": null
	@# }
	@VAULT_ADDR=http://127.0.0.1:8200 VAULT_TOKEN=${LOCAL_VAULT_ROOT_TOKEN} vault read -format=json database/creds/my-role

.PHONY: test-plugin
