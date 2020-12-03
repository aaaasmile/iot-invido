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
	UseTLS         bool
	Influx         *idl.Influx
	SensorCfg      SensorConfig
}

type SensorConfig struct {
	Name string
	ID   string
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
	token, err := readTokenFile(Current.TokenFilename, &Current.SensorCfg)
	if err != nil {
		log.Fatal("Credential error. Please make sure that an account has been initiated. Error is: ", err)
	}
	Current.Influx.Token = token
	log.Printf("Token is %s...", token[:10])
	log.Printf("ID is %s...", Current.SensorCfg.ID[:4])
	log.Printf("Name is %s ", Current.SensorCfg.Name)
	return Current
}

func readTokenFile(fname string, sensCfg *SensorConfig) (string, error) {

	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	cred := struct {
		Token string
		Name  string
		ID    string
	}{}

	err = json.NewDecoder(f).Decode(&cred)
	if err != nil {
		return "", err
	}

	sensCfg.ID = cred.ID
	sensCfg.Name = cred.Name

	return cred.Token, nil
}
