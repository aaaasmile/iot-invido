package conf

import (
	"encoding/json"
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
	TokenFilename  string
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
	token, err := tokenFromFile(Current.TokenFilename)
	if err != nil {
		log.Fatal("Credential error. Please make sure that an account has been initiated. Error is: ", err)
	}
	Current.Influx.Token = token
	log.Printf("Token is %s...", token[:10])
	return Current
}

func tokenFromFile(fname string) (string, error) {

	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	cred := struct {
		Token string
	}{}

	err = json.NewDecoder(f).Decode(&cred)
	if err != nil {
		return "", err
	}
	return cred.Token, nil
}
