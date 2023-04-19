package config

import (
	"os"

	"github.com/Cyan903/c-share/pkg/log"

	"gopkg.in/yaml.v2"
)

type Conf struct {
	DSN         string   `yaml:"DSN,omitempty"`
	Port        int      `yaml:"PORT,omitempty"`
	Mode        string   `yaml:"MODE,omitempty"`
	JWTSecret   string   `yaml:"SECRET,omitempty"`
	UploadLimit int64    `yaml:"UPLOAD_LIMIT,omitempty"`
	UploadPath  string   `yaml:"UPLOAD_PATH,omitempty"`
	CorsAllow   []string `yaml:"CORS-ALLOW,omitempty"`

	Cache struct {
		Address  string `yaml:"ADDRESS,omitempty"`
		Password string `yaml:"Password,omitempty"`
		DB       int    `yaml:"DB,omitempty"`
	} `yaml:"CACHE"`

	Mail struct {
		User     string `yaml:"USER,omitempty"`
		Password string `yaml:"PASSWORD,omitempty"`
		Host     string `yaml:"HOST,omitempty"`
		Port     int    `yaml:"PORT,omitempty"`
	} `yaml:"MAIL"`
}

var Data Conf
var Dev bool

func LoadConfig() Conf {
	config, err := os.ReadFile("config.yaml")
	cfg := Conf{}

	if err != nil {
		log.Error.Println("Could not find config.yaml.")
		os.Exit(1)
	}

	if err = yaml.Unmarshal(config, &cfg); err != nil {
		log.Error.Println("Could not unmarshal config -", err)
		os.Exit(1)
	}

	Data = cfg
	Dev = cfg.Mode == "development"

	return cfg
}
