package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	// Config -.
	Config struct {
		App  `yaml:"app"`
		HTTP `yaml:"http"`
		Log  `yaml:"logger"`
		PG   `yaml:"postgres"`
		//Jwt  `json:"jwt"`
		Redis `yaml:"redis"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
	//Jwt struct {
	//	Secret string `env-required:"true" yaml:"secret_key"   env:"secret_key"`
	//}

	PG struct {
		PoolMax  int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		Username string `env-required:"true" yaml:"username" env-default:"postgres"`
		Password string `env-required:"true" yaml:"password" env-default:"postgres"`
		Host     string `env-required:"true" yaml:"host" env-default:"localhost"`
		Port     string `env-required:"true" yaml:"port" env-default:"5432"`
		Name     string `env-required:"true" yaml:"name" env-default:"postgres"`
	}
	Redis struct {
		Address  string `env-required:"true" yaml:"address" env:"address"`
		Password string `env-required:"true" yaml:"password" env-default:"password"`
		Db       int    `env-required:"true" yaml:"db" env-default:"db"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	currPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := currPath + os.Getenv("path") + "/" + os.Getenv("env") + ".yml"
	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	var cfg *Config

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil

}
