package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sync"
)

type Config struct {
	Database struct {
		Host        string `yaml:"host"`
		Database    string `yaml:"database"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		MaxPoolSize uint64 `yaml:"max_pool_size"`
	} `yaml:"database"`
	Session struct {
		Secret string `yaml:"secret"`
	}
	Social struct {
		Google struct {
			ClientID     string `yaml:"client_id"`
			ClientSecret string `yaml:"client_secret"`
			RedirectURL  string `yaml:"redirect_url"`
		} `yaml:"google"`
	} `yaml:"social"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

var once sync.Once

func NewConfig() (*Config, error) {
	var instance *Config
	once.Do(func() {
		path := os.Getenv("CONFIG_PATH")
		if path == "" {
			log.Println("CONFIG_PATH is not set, using default path /app/config-dev.yaml")
			path = "/app/config-dev.yaml"
		}
		yamlFile, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("error reading config file: %v", err)
		}

		err = yaml.Unmarshal(yamlFile, &instance)
		if err != nil {
			log.Fatalf("error unmarshalling config file: %v", err)
		}
	})
	return instance, nil
}
