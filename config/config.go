package config

type Config struct {
	Version string
	Server  struct {
		Port string
	}
	Redis RedisConfig
}
type RedisConfig struct {
	Host string
	Port string
	DB   int
}
