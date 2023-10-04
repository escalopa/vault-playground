package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	vaultAddr string
	rootToken string

	tlsEnabled bool
	serverCert string
	clientCert string
	clientKey  string
)

func init() {
	flag.StringVar(&vaultAddr, "address", "http://127.0.0.1:8200", "Vault address")
	flag.StringVar(&rootToken, "token", "", "Vault root token")

	flag.BoolVar(&tlsEnabled, "tls", false, "Enable TLS")
	flag.StringVar(&serverCert, "server-cert", "", "Server certificate")
	flag.StringVar(&clientCert, "client-cert", "", "Client certificate")
	flag.StringVar(&clientKey, "client-key", "", "Client key")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	// prepare a client
	client, err := newVaultClient()
	if err != nil {
		log.Fatalf("failed to create vault client: %v", err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(rootToken); err != nil {
		log.Fatalf("failed to set token: %v", err)
	}

	// read the secret
	s, err := client.Secrets.KvV2Read(ctx, "db", vault.WithMountPath("secret"))
	if err != nil {
		log.Fatalf("failed to read secret: %v", err)
	}

	// get the dsn
	dsn, ok := s.Data.Data["dsn"].(string)
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
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
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

func newVaultClient() (*vault.Client, error) {
	opts := []vault.ClientOption{
		vault.WithAddress(vaultAddr),
		vault.WithRequestTimeout(30 * time.Second),
	}

	if tlsEnabled {
		tlsConfig := vault.TLSConfiguration{
			InsecureSkipVerify: false,
			ServerCertificate: vault.ServerCertificateEntry{
				FromFile: serverCert,
			},
			ClientCertificate: vault.ClientCertificateEntry{
				FromFile: clientCert,
			},
			ClientCertificateKey: vault.ClientCertificateKeyEntry{
				FromFile: clientKey,
			},
		}
		opts = append(opts, vault.WithTLS(tlsConfig))
	}

	return vault.New(opts...)
}
