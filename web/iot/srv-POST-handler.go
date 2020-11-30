package iot

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/web/iot/datahandler"
)

func handlePost(w http.ResponseWriter, req *http.Request) error {
	var err error
	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "PubData":
		err = handlePubData(w, req)
	case "FetchData":
		hd := datahandler.HandleData{
			Influx: conf.Current.Influx,
		}
		err = hd.HandleFetchData(w, req)
	case "InsertTestData":
		hd := datahandler.HandleData{
			Influx: conf.Current.Influx,
		}
		err = hd.HandleTestInsertLine(w, req)
	default:
		return fmt.Errorf("%s method is not supported", lastPath)
	}

	return err
}

func handlePubData(w http.ResponseWriter, req *http.Request) error {
	// TODO... parse the data
	return nil
}
