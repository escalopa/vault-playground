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

2. Run the app
```shell
make run
```

## Milestones 🚀

### v1.0.0 🎯
- Read secrets from Vault
- Write secrets to Vault
- Connect to databse with secrets from Vault