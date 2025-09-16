package configs

import (
	"flag"
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configs struct {
	Env string `yaml:"env" env-default:"local"`
	Storage_path string `yaml:"storage_path" env-required:"true"`
	Token_TTL string `yaml:"token_ttl" env-required:"true"`
	JWT_Secret string `yaml:"jwt_secret" env-required:"12345"`
}

func MustLoad() *Configs {
	const op = "AuthService.internal.configs.Mustload"
	
	path := fetchConfigPath()
	if path == "" {
		panic("Invalid path. Erorr: " + op)
	}

	var cfg Configs

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Invalid read config. Erorr: " + op)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	fmt.Println(res)

	if res == "" {
		res = os.Getenv("CONFIG")
	}

	return res
}