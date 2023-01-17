package config

type System struct {
	Protocol string
	Socket   string
}

func (s System) Set(cfg *Configs) error {
	cfg.System = s
	return nil
}

func WithSystem(system System) Config {
	return system
}
