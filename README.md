# vault-playground 🔒
An enviroment to play with Hashicorp Vault &amp; learn it in depth

## Pre-requisites 📚

- [Docker](https://www.docker.com/)
- [Vault](https://www.vaultproject.io/)
- [Go](https://go.dev/doc/install)

## Run 💨

1. Start the docker containers
```shell
docker compose up -d
```

2. Run the migration scripts
```shell
make migrate
```

3. Start vault in insecure mode
```shell
make vault
```

4. Set vault address in shell (On TLS use `https`)
```shell
export VAULT_ADDR=http://127.0.0.1:8200
```

5. Enable role path in vault
```shell
vault auth enable approle
```

6. Create role for app
```shell
make role-create
```

7. Set role policy
```shell
make policy-create
```

8. Set database dsn secret
```shell
make dns-create
```

9. Run the app
```shell
make run
```

10. Get user orders
```shell
curl -q http://localhost:8080/order/101 | jq
```

--- 

## Production Use Case 🏘

1. Create directory `./vault/data`
```shell
mkdir -p ./vault/data
```

2. Start server 
```shell 
sudo make vault-prod
```

3. Init the server
```shell
vault operator init
```

4. Unseal the server using 3 secrets, Secrets can be found in the output of the 3rd command
```shell
vault operator unseal
```

5. Login to the server, Token can be found in the output of the 3rd command
```shell
vault login
```

## Milestones 🚀

### v1.0.0 🎯
- Read secrets from Vault
- Write secrets to Vault
- Connect to databse with secrets from Vault

### v2.0.0 🎯
- Create approle, policy for application
- Create a vault client with approle
- Read secrets from vault with approle policies
- Run vault in production mode

### v3.0.0 🎯
- 