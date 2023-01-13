package config

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
