package config

import "github.com/spf13/cobra"

type Auth struct {
	CaCert     string
	ClientCert string
	ClientKey  string
	ServerName string
}

func (a Auth) Set(cfg *Configs) error {
	cfg.Auth = a
	return nil
}

func WithAuth(auth Auth) Config {
	return auth
}

func (a *Auth) AddFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&a.CaCert, "ca_crt", "/etc/aurae/ca.crt", "The CA certificate")
	cmd.Flags().StringVar(&a.ClientCert, "client_crt", "/etc/aurae/_signed.client.crt", "The client certificate")
	cmd.Flags().StringVar(&a.ClientKey, "client_key", "/etc/aurae/client.key", "The client certificate key")
}
