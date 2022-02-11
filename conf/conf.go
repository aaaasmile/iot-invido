package conf

import (
	"encoding/json"
	"io/ioutil"
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
	SQLite         SQLite
}

type SQLite struct {
	DBPath string
}

type SensorConfig struct {
	Name string
	ID   string
}

type TokenFile struct {
	Token string
	Name  string
	ID    string
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

func CreateEmptyTokenFile() error {
	fname := "token.json"
	cred := TokenFile{
		ID:    "enter your ID for influxdb",
		Token: "enter your Token for influxdb",
		Name:  "enter your Name for influxdb",
	}

	raw, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fname, raw, 0664)
	if err != nil {
		return err
	}
	log.Println("Toke file created: ", fname)
	return nil
}

func readTokenFile(fname string, sensCfg *SensorConfig) (string, error) {

	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()
	cred := TokenFile{}

	err = json.NewDecoder(f).Decode(&cred)
	if err != nil {
		return "", err
	}

	sensCfg.ID = cred.ID
	sensCfg.Name = cred.Name

	return cred.Token, nil
}
