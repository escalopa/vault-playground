# Vault Cofig
VAULT_ADDR=http://127.0.0.1:8200
VAULT_TOKEN=my-root-token
TLS=false

# Application Config
APP_NAME=vault-playground
PORT=8080

# Application Vault Config
ROLE_ID=$(APP_NAME)-role
POLICY_ID=$(APP_NAME)-policy

run:
	APP_ROLE_ID="$$(vault read -field=role_id auth/approle/role/$(ROLE_ID)/role-id )"; \
	APP_SECRET_ID="$$(vault write  -f -field=secret_id auth/approle/role/$(ROLE_ID)/secret-id )"; \
	go run cmd/*.go \
		--app-name=$(APP_NAME) \
		--port=$(PORT) \
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

vault-prod:
	vault server -config=config.hcl

vault-db:
	POSTGRES_URL="localhost:5432"; \
	vault write database/config/postgresql \
		plugin_name=postgresql-database-plugin \
		connection_url="postgresql://{{username}}:{{password}}@$${POSTGRES_URL}/postgres?sslmode=disable" \
		allowed_roles=readonly \
		username="postgres" \
		password="postgres" 

role-create:
	vault write auth/approle/role/$(ROLE_ID) \
    secret_id_ttl=10m \
    token_num_uses=10 \
    token_ttl=20m \
    token_max_ttl=30m \
    secret_id_num_uses=40 \
    token_policies=$(POLICY_ID)

policy-create:
	vault policy write $(POLICY_ID) ./policy.hcl

db-role-create:
	vault write database/roles/readonly \
      db_name=postgresql \
      creation_statements=@readonly.sql \
      default_ttl=1h \
      max_ttl=24h

dsn-create:
	vault kv put --mount=secret $(APP_NAME)/db dsn="postgres://postgres:postgres@localhost:5432/vault-db"

migrate:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/vault-db?sslmode=disable" up