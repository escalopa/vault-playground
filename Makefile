VAULT_ADDR=http://127.0.0.1:8200
VAULT_TOKEN=my-root-token
TLS=false

run:
	go run main.go \
		--address=$(VAULT_ADDR) \
		--token=$(VAULT_TOKEN) \
		--tls=$(TLS) \
		--server-cert="$(PWD)/cert/vault-ca.pem" \
		--client-cert="$(PWD)/cert/vault-cert.pem" \
		--client-key="$(PWD)/cert/vault-key.pem" 
		
vault:
	vault server -dev -dev-root-token-id=$(VAULT_TOKEN)

vault-tls:
	mkdir -p cert
	vault server -dev -dev-tls -dev-tls-cert-dir="$PWD/cert" -dev-root-token-id=$(VAULT_TOKEN)
