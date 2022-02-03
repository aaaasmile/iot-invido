package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/web"
	"github.com/aaaasmile/iot-invido/web/idl"
)

func main() {
	var ver = flag.Bool("ver", false, "Prints the current version")
	var configfile = flag.String("config", "config.toml", "Configuration file path")
	var createtkfile = flag.Bool("createtkfile", false, "Create an empty token file")
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

	if err := web.RunService(*configfile); err != nil {
		panic(err)
	}
}
