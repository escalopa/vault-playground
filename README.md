# vault-playground 🔒
An enviroment to play with Hashicorp Vault &amp; learn it in depth

## Pre-requisites 📚

- [Docker](https://www.docker.com/)
- [Vault](https://www.vaultproject.io/)
- [Go](https://go.dev/doc/install)

## Run 💨

1. Start the docker containers (PG/Redis)
```shell
docker compose up -d
```

2. Start vault in insecure mode
```shell
make vault
```

3. Set vault address in shell (On TLS use `https`)
```shell
export VAULT_ADDR=http://127.0.0.1:8200
```

4. Enable role path in vault
```shell
vault auth enable approle
```

5. Create role for app
```shell
make role-create
```

6. Set role policy
```shell
make policy-create
```

7. Set database dsn secret
```shell
make dns-create
```

8. Run the app
```shell
make run
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