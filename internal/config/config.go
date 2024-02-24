package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

const (
	EnvLocal string = "local"
	EnvDev   string = "dev"
	EnvProd  string = "prod"

	Secret = "SECRET"

	PostgresUser     = "POSTGRES_USER"
	PostgresPassword = "POSTGRES_PASSWORD"
	PostgresPort     = "POSTGRES_PORT"
	PostgresDB       = "POSTGRES_DB"
	PostgresHost     = "POSTGRES_HOST"

	RedisHost          = "REDIS_HOST"
	RedisPort          = "REDIS_PORT"
	RedisPassword      = "REDIS_PASSWORD"
	RedisDatabaseIndex = "REDIS_DB_INDEX"
)

type (
	Config struct {
		Application ApplicationConfig `yaml:"application" env-required:"true"`
		Handler     HandlerConfig     `yaml:"handler" env-required:"true"`
		Token       TokenConfig       `yaml:"token" env-required:"true"`
		Server      ServerConfig      `yaml:"server" env-required:"true"`
	}

	ApplicationConfig struct {
		Env   string `yaml:"env" env-default:"local"`
		Debug bool   `yaml:"debug" env-default:"true"`
	}

	HandlerConfig struct {
		ContextTimeout time.Duration `yaml:"context-timeout" env-required:"true"`
	}

	ServerConfig struct {
		Port         int           `yaml:"port" env-default:"8091"`
		Host         string        `yaml:"host" env-required:"true"`
		WriteTimeout time.Duration `yaml:"write-timeout" env-required:"true"`
		ReadTimeout  time.Duration `yaml:"read-timeout" env-required:"true"`
	}

	TokenConfig struct {
		AccessTTL  time.Duration `yaml:"access_ttl" env-required:"true"`
		RefreshTTL time.Duration `yaml:"refresh_ttl" env-required:"true"`
	}
)

func MustLoad() *Config {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Failed donwload the file .env: %v", err))
	}

	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config
	err = cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config: " + err.Error())
	}

	switch cfg.Application.Env {
	case EnvLocal, EnvDev, EnvProd:
		fmt.Println("Env: ", cfg.Application.Env)
	default:
		panic("the env is not specified correctly: " + cfg.Application.Env)
	}
	checkEnv()
	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}

func checkEnv() {
	var envKeys = []string{
		Secret,
	}
	for _, key := range envKeys {
		v := os.Getenv(key)
		if v == "" {
			panic("the value for the key " + key + " is not set")
		}
	}
}
