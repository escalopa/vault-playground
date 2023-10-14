# vault-playground 🔒
An enviroment to play with Hashicorp Vault &amp; learn it in depth

## Pre-requisites 📚

- [Docker](https://www.docker.com/)
- [Vault](https://www.vaultproject.io/)
- [Go](https://go.dev/doc/install)

## Run 💨

1. Start the docker containers
```bash
docker compose up -d
```

2. Run the migration scripts
```bash
make migrate
```

3. Start vault in insecure mode
```bash
make vault
```

4. Set vault address in shell (On TLS use `https`)
```bash
export VAULT_ADDR=http://127.0.0.1:8200
```

5. Enable role path in vault
```bash
vault auth enable approle
```

6. Create role for app
```bash
make role-create
```

7. Set role policy
```bash
make policy-create
```

8. Set database dsn secret
```bash
make dns-create
```

9. Run the app
```bash
make run
```

10. Get user orders
```bash
curl -q http://localhost:8080/order/101 | jq
```

--- 

## Production Use Case 🏘

1. Create directory `./vault/data`
```bash
mkdir -p ./vault/data
```

2. Start server 
```shell 
sudo make vault-prod
```

3. Init the server
```bash
vault operator init
```

4. Unseal the server using 3 secrets, Secrets can be found in the output of the 3rd command
```bash
vault operator unseal
```

5. Login to the server, Token can be found in the output of the 3rd command
```bash
vault login
```

## Database & Dynamic Secrets 🗄

Before we start make sure vault is up and running

1. Create a role for vault in postgres db
```bash
docker exec -i  db psql -U postgres -c "CREATE ROLE \"vault-ro\" NOINHERIT;"
```

2. Grant the ability to read all tables to vault role
```bash
docker exec -i  db psql -U postgres -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO \"vault-ro\";"
``` 

3. Enable database secrets engine
```bash
vault secrets enable database
```

4. Create database configuration in vault
```bash
make vault-db
```

5. Create database role in vault
```bash
make db-role-create
```

6. Get sample database credentials
```bash
vault read database/creds/readonly
```

7. Check the database credentials in postgres
```bash
docker exec -i \       
    db \
    psql -U postgres -c "SELECT usename, valuntil FROM pg_user;"
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
- Link vault with postgres
- Create/Read dynamic secrets