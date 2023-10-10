VAULT_ADDR=http://127.0.0.1:8200
VAULT_TOKEN=my-root-token
TLS=false

APP_NAME=vault-playground
ROLE_ID=application-role
POLOLICY_ID=application-policy

run:
	APP_ROLE_ID="$$(vault read -field=role_id auth/approle/role/$(ROLE_ID)/role-id )"; \
	APP_SECRET_ID="$$(vault write  -f -field=secret_id auth/approle/role/$(ROLE_ID)/secret-id )"; \
	go run main.go \
		--app-name=$(APP_NAME) \
		--address=$(VAULT_ADDR) \
		--role-id=$$APP_ROLE_ID \
		--secret-id=$$APP_SECRET_ID \
		--tls=$(TLS) \
		--server-cert="$(PWD)/cert/vault-ca.pem" \
		--client-cert="$(PWD)/cert/vault-cert.pem" \
		--client-key="$(PWD)/cert/vault-key.pem"

vault:
	vault server -dev -dev-root-token-id=$(VAULT_TOKEN)

vault-tls:
	mkdir -p cert
	vault server -dev -dev-tls -dev-tls-cert-dir="$(PWD)/cert" -dev-root-token-id=$(VAULT_TOKEN)

path-create:
	vault secrets enable -path=application/$(APP_NAME) kv

role-create:
	vault write auth/approle/role/$(ROLE_ID) \
    secret_id_ttl=10m \
    token_num_uses=10 \
    token_ttl=20m \
    token_max_ttl=30m \
    secret_id_num_uses=40

role-read:
	vault read -field=role_id auth/approle/role/$(ROLE_ID)/role-id

role-secret:
	vault write -f -field=secret_id auth/approle/role/$(ROLE_ID)/secret-id

role-policy:
	vault write auth/approle/role/$(ROLE_ID) policies=$(POLOLICY_ID)

db-dsn-create:
	vault kv put application/$(APP_NAME)/db dsn="postgres://postgres:postgres@localhost:5432/vault-db"