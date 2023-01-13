package config

import (
	"fmt"
	"os/user"
)

type Configs struct {
	Auth   Auth
	System System
}

type Config interface {
	Set(p *Configs) error
}

func Default() (*Configs, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Configs{
		Auth: Auth{
			CaCert:     fmt.Sprintf("%s/.aurae/pki/ca.crt", usr.HomeDir),
			ClientCert: fmt.Sprintf("%s/.aurae/pki/_signed.client.nova.crt", usr.HomeDir),
			ClientKey:  fmt.Sprintf("%s/.aurae/pki/client.nova.key", usr.HomeDir),
			ServerName: "server.unsafe.aurae.io",
		},
		System: System{
			Protocol: "unix",
			Socket:   "/var/run/aurae/aurae.sock",
		},
	}, nil
}

func From(cfg ...Config) (*Configs, error) {
	c, err := Default()
	if err != nil {
		return nil, err
	}

	for _, config := range cfg {
		err := config.Set(c)

		if err != nil {
			return nil, err
		}
	}

	return c, nil
}
