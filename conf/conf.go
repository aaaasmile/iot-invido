package conf

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/aaaasmile/iot-invido/web/idl"
)

type Config struct {
	ServiceURL     string
	RootURLPattern string
	DebugVerbose   bool
	VueLibName     string
	Influx         *idl.Influx
}

var Current = &Config{}

func ReadConfig(configfile string) *Config {
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := toml.DecodeFile(configfile, &Current); err != nil {
		log.Fatal(err)
	}
	return Current
}
