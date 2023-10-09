VAULT_ADDR=https://127.0.0.1:8200
VAULT_TOKEN=my-root-token
TLS=false

ROLE_ID=fa87f98a-fb38-4919-abbc-4d0b289f08a4
SECRET_ID=93ad23c1-94bb-483d-b8d7-21a51c0389f7

run:
	go run main.go \
		--address=$(VAULT_ADDR) \
		--role-id=$(ROLE_ID) \
		--secret-id=$(SECRET_ID) \
		--tls=$(TLS) \
		--server-cert="$(PWD)/cert/vault-ca.pem" \
		--client-cert="$(PWD)/cert/vault-cert.pem" \
		--client-key="$(PWD)/cert/vault-key.pem" 
		
vault:
	vault server -dev -dev-root-token-id=$(VAULT_TOKEN)

vault-tls:
	mkdir -p cert
	vault server -dev -dev-tls -dev-tls-cert-dir="$(PWD)/cert" -dev-root-token-id=$(VAULT_TOKEN)

role-create:
	vault write auth/approle/login \
		role_id=$(ROLE_ID) \
		secret_id=$(SECRET_ID)
