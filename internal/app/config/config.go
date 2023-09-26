package config

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type MainConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	Redis    CacheConfig
	JWT      JWTConfig
	Log      log.FieldLogger
}

type ServerConfig struct {
	Address int
	RunMode string
}

type JWTConfig struct {
	Secret                     string
	RefreshSecret              string
	AccessTokenExpireDuration  int
	refreshTokenExpireDuration int
}

type DatabaseConfig struct {
	Host   string
	Port   string
	DBName string
	User   string
	Pass   string
}

type CacheConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	UseTLS   bool
}

func GetConfig() (*MainConfig, error) {
	cfgPath := getConfigPath(os.Getenv("APP_ENV"))
	v, err := LoadConfig(cfgPath, "yml")
	if err != nil {
		log.Errorf("Error in load config %v", err)
		return nil, err
	}

	cfg, err := ParseConfig(v)
	if err != nil {
		log.Errorf("Error in parse config %v", err)
		return nil, err
	}
	return cfg, err
}

func ParseConfig(v *viper.Viper) (*MainConfig, error) {
	var cfg MainConfig
	err := v.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return nil, err
	}
	return &cfg, nil
}

func LoadConfig(filename string, fileType string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigType(fileType)
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Unable to read config: %v", err)
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}
	return v, nil
}

func getConfigPath(env string) string {
	if env == "docker" {
		return "/app/config/config-docker"
	} else {
		return "configs/config-development"
	}
}
