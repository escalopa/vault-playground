package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
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
		log.Fatal(err)
	}

	// authenticate with a root token (insecure)
	if err := client.SetToken(rootToken); err != nil {
		log.Fatal(err)
	}

	// write a secret
	_, err = client.Secrets.KvV2Write(ctx, "foo", schema.KvV2WriteRequest{
		Data: map[string]any{
			"password1": "abc123",
			"password2": "correct horse battery staple",
		}},
		vault.WithMountPath("secret"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret written successfully")

	// read the secret
	s, err := client.Secrets.KvV2Read(ctx, "foo", vault.WithMountPath("secret"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("secret retrieved:", s.Data.Data)
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
