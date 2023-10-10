package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	vault "github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	appName string

	vaultAddr string
	rootToken string

	roleID   string
	secretID string

	tlsEnabled bool
	serverCert string
	clientCert string
	clientKey  string
)

func init() {
	flag.StringVar(&appName, "app-name", "vault-playground", "App name")

	flag.StringVar(&vaultAddr, "address", "http://127.0.0.1:8200", "Vault address")
	flag.StringVar(&rootToken, "token", "", "Vault root token")

	flag.StringVar(&roleID, "role-id", "", "Vault role ID")
	flag.StringVar(&secretID, "secret-id", "", "Vault secret ID")

	flag.BoolVar(&tlsEnabled, "tls", false, "Enable TLS")
	flag.StringVar(&serverCert, "server-cert", "", "Server certificate")
	flag.StringVar(&clientCert, "client-cert", "", "Client certificate")
	flag.StringVar(&clientKey, "client-key", "", "Client key")

	flag.Parse()
}

func main() {
	ctx := context.Background()

	// prepare a client
	client, err := newVaultClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// read the secret
	homePath := fmt.Sprintf("application/%s", appName)
	s, err := client.KVv2(homePath).Get(ctx, "db")
	if err != nil {
		log.Fatalf("failed to read secret: %v", err)
	}

	// get the dsn
	dsn, ok := s.Data["dsn"].(string)
	if !ok {
		log.Fatal("failed to get dsn")
	}

	// connect to the database
	pool, err := newPostgresClient(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// close the connection
	defer pool.Close()

	// simple query
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("failed to acquire connection: %v", err)
	}

	defer conn.Release()

	_, err = conn.Exec(ctx, "SELECT 1")
	if err != nil {
		log.Fatalf("failed to execute query: %v", err)
	}

	log.Println("Done")
}

func newPostgresClient(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

func newVaultClient(ctx context.Context) (*vault.Client, error) {
	config := vault.DefaultConfig()
	config.Address = vaultAddr

	if tlsEnabled {
		config.ConfigureTLS(&vault.TLSConfig{
			ClientCert: clientCert,
			ClientKey:  clientKey,
			CACert:     serverCert,
			Insecure:   false,
		})
	}

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %v", err)
	}

	appRoleAuth, err := auth.NewAppRoleAuth(roleID, &auth.SecretID{FromString: secretID})
	if err != nil {
		return nil, fmt.Errorf("failed to create app role: %v", err)
	}

	authInfo, err := client.Auth().Login(ctx, appRoleAuth)
	if err != nil {
		return nil, fmt.Errorf("failed to login using app role: %v", err)
	}
	if authInfo == nil {
		return nil, fmt.Errorf("no auth info was returned after login")
	}

	client.SetToken(authInfo.Auth.ClientToken)

	return client, nil
}
