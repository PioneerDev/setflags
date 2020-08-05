package setting

import (
	"io/ioutil"
	"log"

	"github.com/gofrs/uuid"
	yaml "gopkg.in/yaml.v2"
)

/*
Config is used in user system.
*/
type Config struct {
	RUNMODE     string  `yaml:"run_mode"`
	RewardRatio float64 `yaml:"reward_ratio"`

	App struct {
		Name         string `yaml:"name"`
		PageSize     int    `yaml:"page_size"`
		HTTPPort     int    `yaml:"http_port"`
		ReadTimeOut  int    `yaml:"read_timeout"`
		WriteTimeOut int    `yaml:"write_timeout"`
		JWTSecret    string `yaml:"jwt_secret"`
		// hour
		JWTTokenExpireTime int `yaml:"jwt_token_expire_time"`
	}

	DataBase struct {
		Type     string `yaml:"type"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	Bot struct {
		ClientID     uuid.UUID `yaml:"client_id"`
		ClientSecret string    `yaml:"client_secret"`
		SessionID    string    `yaml:"session_id"`
		Pin          string    `yaml:"pin"`
		PinToken     string    `yaml:"pin_token"`
		PrivateKey   string    `yaml:"private_key"`
		CodeVerifier string    `yaml:"code_verifier"`
	}
}

var cfg *Config

// LoadConfig load config
func LoadConfig() (*Config, error) {
	cfg = new(Config)
	bytes, err := ioutil.ReadFile("/api/secrets/config.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal([]byte(bytes), cfg)
	if err != nil {
		panic(err)
	}
	return cfg, nil
}

// GetConfig get config
func GetConfig() *Config {
	var err error
	if cfg == nil {
		if cfg, err = LoadConfig(); err != nil {
			log.Panicf("Failed to load config: %s\n", err)
		}
	}
	return cfg
}
