package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/db/sqlite"
	"github.com/aaaasmile/iot-invido/web"
	"github.com/aaaasmile/iot-invido/web/idl"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var createtkfile = flag.Bool("createtkfile", false, "Create an empty token file")
	var createnewuser = flag.Bool("createuser", false, "Create a new user in the sqlite db")
	flag.Parse()

	if *ver {
		fmt.Printf("%s, version: %s", idl.Appname, idl.Buildnr)
		os.Exit(0)
	}
	if *createtkfile {
		if err := conf.CreateEmptyTokenFile(); err != nil {
			fmt.Println("Fatal error: ", err)
		}
		os.Exit(0)
	}
	if *createnewuser {
		if err := sqlite.CreateNewUser(*configfile); err != nil {
			fmt.Println("Fatal error: ", err)
		}
		log.Println("That's all for crate new user")
		os.Exit(0)
	}

	if err := web.RunService(*configfile); err != nil {
		panic(err)
	}
}
