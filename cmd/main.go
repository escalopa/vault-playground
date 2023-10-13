package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/escalopa/vault-playground/internal/handler"
	"github.com/escalopa/vault-playground/internal/pg"
	"github.com/escalopa/vault-playground/internal/service"
)

var (
	appName string
	port    string

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
	flag.StringVar(&port, "port", "8080", "App port")

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

	// get the database dsn
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

	storageFacade := pg.New(pool)
	ordersService := service.NewOrdersService(storageFacade)

	h := handler.New(ordersService)

	log.Default().Printf("starting server at 0.0.0.0:%s", port)
	err = h.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
