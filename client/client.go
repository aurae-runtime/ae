/* -------------------------------------------------------------------------- *\
 *             Apache 2.0 License Copyright © 2022 The Aurae Authors          *
 *                                                                            *
 *                +--------------------------------------------+              *
 *                |   █████╗ ██╗   ██╗██████╗  █████╗ ███████╗ |              *
 *                |  ██╔══██╗██║   ██║██╔══██╗██╔══██╗██╔════╝ |              *
 *                |  ███████║██║   ██║██████╔╝███████║█████╗   |              *
 *                |  ██╔══██║██║   ██║██╔══██╗██╔══██║██╔══╝   |              *
 *                |  ██║  ██║╚██████╔╝██║  ██║██║  ██║███████╗ |              *
 *                |  ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝╚══════╝ |              *
 *                +--------------------------------------------+              *
 *                                                                            *
 *                         Distributed Systems Runtime                        *
 *                                                                            *
 * -------------------------------------------------------------------------- *
 *                                                                            *
 *   Licensed under the Apache License, Version 2.0 (the "License");          *
 *   you may not use this file except in compliance with the License.         *
 *   You may obtain a copy of the License at                                  *
 *                                                                            *
 *       http://www.apache.org/licenses/LICENSE-2.0                           *
 *                                                                            *
 *   Unless required by applicable law or agreed to in writing, software      *
 *   distributed under the License is distributed on an "AS IS" BASIS,        *
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. *
 *   See the License for the specific language governing permissions and      *
 *   limitations under the License.                                           *
 *                                                                            *
\* -------------------------------------------------------------------------- */

package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/aurae-runtime/ae/client/pkg/config"
	"github.com/aurae-runtime/ae/client/pkg/runtime"
)

var ErrOptionNotAvailable = errors.New("option is not available")

type AuraeClient interface {
	Runtime() (runtime.Runtime, error)
}

type auraeClient struct {
	cfg     *config.Configs
	conn    grpc.ClientConnInterface
	runtime runtime.Runtime
}

func New(ctx context.Context, cfg ...config.Config) (AuraeClient, error) {
	cf, err := config.From(cfg...)
	if err != nil {
		log.Fatal("Cannot initialize config", err)
	}

	tlsCredentials, err := loadTLSCredentials(cf.Auth)
	if err != nil {
		log.Fatal("Cannot load TLS credentials", err)
	}

	dialer := func(ctx context.Context, addr string) (net.Conn, error) {
		d := net.Dialer{}

		return d.DialContext(ctx, cf.System.Protocol, addr)
	}

	conn, err := grpc.Dial(
		cf.System.Socket,
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithContextDialer(dialer),
	)
	if err != nil {
		log.Fatal("Cannot Dial", err)

		return nil, err
	}

	r, err := runtime.New(ctx, conn)
	if err != nil {
		log.Fatal("Cannot crete runtime client", err)

		return nil, err
	}

	c := &auraeClient{
		cfg:     cf,
		conn:    conn,
		runtime: r,
	}

	return c, nil
}

func loadTLSCredentials(auth config.Auth) (credentials.TransportCredentials, error) {
	caPEM, err := os.ReadFile(auth.CaCert)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caPEM) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
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

func (c *auraeClient) Runtime() (runtime.Runtime, error) {
	if c.runtime == nil {
		return nil, fmt.Errorf("configuration: %w", ErrOptionNotAvailable)
	}

	return c.runtime, nil
}
