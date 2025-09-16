package configs

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Configs struct {
	Env string `yaml:"env" env-default:"local"`
	MongoUri string `yaml:"mongo_uri" env-default:"none"`
	Raddr string `yaml:"redis_addr" env-required:"true" env-default:"none"`
	Rpassword string `yaml:"redis_password" env-required:"true"`
	RabbitAddr string `yaml:"rabbit_addr" env-required:"true" env-default:"none"`
	HTTPServer `yaml:"http_server"`
	KafkaBroker `yaml:"kafka"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type KafkaBroker struct {
	Broker string `yaml:"broker" env-required:"true" env-default:"localhost:29092"`
	Topic string `yaml:"topic" env-default:"test-topic"`
	GroupID string `yaml:"groupID" env-default:"1"`
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