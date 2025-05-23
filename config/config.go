package config

var conf *Config

const (
	env     = "local"
	address = "localhost"
	port    = "8080"
)

type Config struct {
	env     string
	address string
	port    string
}

func InitConfig() {
	conf = &Config{
		env:     env,
		address: address,
		port:    port,
	}
}

func GetAddress() string {
	return conf.address
}

func GetPort() string {
	return conf.port
}

func GetEnv() string {
	return conf.env
}
