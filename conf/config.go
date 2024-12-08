package conf

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var cfg *Config

type Config struct {
	DB    DBConfig          `yaml:"database"`
	Redis DBConfig          `yaml:"redis"`
	Web   ApplicationConfig `yaml:"web"`
}

type DBConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Name string `yaml:"name"`
}

type RedisConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type ApplicationConfig struct {
	SaltKey string `yaml:"salt_key"`
	Port    int    `yaml:"port"`
}

func LoadYaml(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func init() {
	var err error
	if cfg, err = LoadYaml("conf/config.yaml"); err != nil {
		log.Fatal(err)
	}
}

func GetDatabaseConfig() *DBConfig {
	return &cfg.DB
}

func GetApplicationConfig() *ApplicationConfig {
	return &cfg.Web
}
