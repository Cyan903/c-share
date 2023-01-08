package config

import (
	"io/ioutil"
	"os"

	"github.com/Cyan903/c-share/pkg/log"

	"gopkg.in/yaml.v2"
)

type conf struct {
	DSN  string `yaml:"DSN,omitempty"`
	Port int    `yaml:"PORT,omitempty"`

	Mode      string `yaml:"MODE,omitempty"`
	JWTSecret string `yaml:"SECRET,omitempty"`
}

var Dev bool

func LoadConfig() conf {
	config, err := ioutil.ReadFile("config.yaml")
	cfg := conf{}

	if err != nil {
		log.Error.Println("Could not find config.yaml.")
		os.Exit(1)
	}

	if err = yaml.Unmarshal(config, &cfg); err != nil {
		log.Error.Println("Could not unmarshal config -", err)
		os.Exit(1)
	}

	return cfg
}
