package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/aurae-runtime/ae/pkg/config"
	"github.com/aurae-runtime/ae/pkg/discovery"
	"github.com/aurae-runtime/ae/pkg/health"
	"github.com/aurae-runtime/ae/pkg/observe"
)

type Client interface {
	Discovery() (discovery.Discovery, error)
	Health() (health.Health, error)
	Observe() (observe.Observe, error)
}

type client struct {
	cfg       *config.Configs
	conn      grpc.ClientConnInterface
	discovery discovery.Discovery
	health    health.Health
	observe   observe.Observe
}

func New(ctx context.Context, cfg ...config.Config) (Client, error) {
	cf, err := config.From(cfg...)
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize client config: %s", err)
	}

	tlsCreds, err := loadTLSCredentials(cf.Auth)
	if err != nil {
		return nil, fmt.Errorf("Failed to load TLS credentials: %s", err)
	}

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(ctx, cf.System.Protocol, addr)
	}

	conn, err := grpc.Dial(cf.System.Socket, grpc.WithTransportCredentials(tlsCreds), grpc.WithContextDialer(dialer))
	if err != nil {
		return nil, fmt.Errorf("Failed to dial server: %s", err)
	}

	return &client{
		cfg:       cf,
		conn:      conn,
		discovery: discovery.New(ctx, conn),
		health:    health.New(ctx, conn),
		observe:   observe.New(ctx, conn),
	}, nil
}

func loadTLSCredentials(auth config.Auth) (credentials.TransportCredentials, error) {
	caPEM, err := os.ReadFile(auth.CaCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("Failed to add server CA's certificate")
	}

	clientKeyPair, err := tls.LoadX509KeyPair(auth.ClientCert, auth.ClientKey)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientKeyPair},
		RootCAs:      certPool,
		ServerName:   auth.ServerName,
	}

	return credentials.NewTLS(config), nil
}

func (c *client) Discovery() (discovery.Discovery, error) {
	if c.discovery == nil {
		return nil, fmt.Errorf("Discovery service is not available")
	}
	return c.discovery, nil
}

func (c *client) Health() (health.Health, error) {
	if c.health == nil {
		return nil, fmt.Errorf("Health service is not available")
	}
	return c.health, nil
}

func (c *client) Observe() (observe.Observe, error) {
	if c.observe == nil {
		return nil, fmt.Errorf("Observe service is not available")
	}
	return c.observe, nil
}
