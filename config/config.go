package config

type Config struct {
	DB     DBConfig
	Server ServerConfig
}

type DBConfig struct {
	Conn  string `toml:"conn"`
	Debug bool   `toml:"debug"`
}

type ServerConfig struct {
	Addr     string `toml:"addr"`
	Debug    bool   `toml:"debug"`
	UseTLS   bool   `toml:"useTLS"` //test only,not necessary in design
	CertFile string `toml:"cert"`
	KeyFile  string `toml:"key"`
}
