package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
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
	s, err := client.KVv2("secret").Get(ctx, fmt.Sprintf("%s/db", appName))
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
