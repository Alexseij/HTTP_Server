package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	DB     *DBConfig
	Server *SeverConfig
}

type SeverConfig struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type DBConfig struct {
	DBUser     string `yaml:"db_user"`
	DBPassword string `yaml:"db_password"`
	DBName     string `yaml:"db_name"`
	DBHost     string `yaml:"db_host"`
}

func LoadCfg(path string) (*Config, error) {
	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
