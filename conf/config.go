package conf

import (
	"crypto/rsa"
	"gitlab/live/be-live-api/model"
	"gitlab/live/be-live-api/service"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
)

var cfg *Config

type Config struct {
	DB    DBConfig `yaml:"database"`
	Redis DBConfig `yaml:"redis"`
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

func SeedRoles(roleService *service.RoleService) {
	roles := []model.Role{
		{Type: string(model.ADMINROLE), Description: "Administrator role"},
		{Type: string(model.USERROLE), Description: "Default user role"},
		{Type: string(model.GUESTROLE), Description: "Guest user role"},
	}

	for _, role := range roles {
		existingRole, _ := roleService.GetRoleByType(role.Type)
		if existingRole != nil {
			continue // Role already exists
		}
		if err := roleService.CreateRole(&role); err != nil {
			log.Fatalf("Failed to seed role: %v", err)
		}
	}

	log.Println("Roles seeded successfully")
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
