# Determine this makefile's path.
# Be sure to place this BEFORE `include` directives, if any.
THIS_FILE := $(lastword $(MAKEFILE_LIST))
MOCKDBSHASUM := $(shell shasum vault/plugins/github-dk | cut -d' ' -f1)

all: export

export:
export VAULT_ADDR=http://127.0.0.1:8200

default: dev

dev-flow: dev disable enable create-role get-creds

start-vault:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./vault/plugins
	echo "Open a new terminal, and run export VAULT_ADDR=http://127.0.0.1:8200"

dev:
	go build -o vault/plugins/github-dk plugin/main.go

disable:
	vault secrets disable database

enable:
	vault secrets enable database
	vault write sys/plugins/catalog/database/github-dk \
    sha_256=$(MOCKDBSHASUM) \
    command="github-dk"
	vault write database/config/github-dk \
     plugin_name="github-dk" \
     url="https://github.dev.us-east-1.aws.galleon.c.statestr.com/api/v3/" \
     apitoken=abc123

create-role:
	vault write database/roles/my-role \
  	db_name=mockdb \
  	creation_statements="CREATE USER \"{{username}}\" WITH PASSWORD '{{password}}'; \
  	     GRANT ALL ON \"vault\" TO \"{{username}}\";" \
  	default_ttl="1h" \
  	max_ttl="24h"

get-creds:
	vault read database/creds/my-role

clean:
	rm -f ./vault/plugins/github-dk
