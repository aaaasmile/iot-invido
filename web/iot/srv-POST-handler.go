package iot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/util"
	"github.com/aaaasmile/iot-invido/web/iot/datahandler"
)

func handlePost(w http.ResponseWriter, req *http.Request) error {
	var err error
	lastPath := getURLForRoute(req.RequestURI)
	log.Println("Check the last path ", lastPath)
	switch lastPath {
	case "PubData":
		hd := datahandler.HandleData{
			Influx:    conf.Current.Influx,
			SensorCfg: &conf.Current.SensorCfg,
		}
		err = hd.HandlePubData(w, req)
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
	case "CheckAPIToken":
		err = handleCheckAPIToken(w, req)
	case "SignIn":
		err = handleSignIn(w, req)
	default:
		return fmt.Errorf("%s method is not supported", lastPath)
	}

	return err
}

func handleCheckAPIToken(w http.ResponseWriter, req *http.Request) error {
	log.Println("Check API Token")
	isvalid, err := validateAPIHeader(req)
	if err != nil {
		return err
	}
	rspdata := struct {
		Valid bool
	}{
		Valid: isvalid,
	}

	return util.WriteJsonResp(w, rspdata)
}

func validateAPIHeader(req *http.Request) (bool, error) {
	tk := req.Header.Get("x-api-sessiontoken")
	if tk == "" {
		return false, nil
	}
	// TODO
	return false, nil
}

func handleSignIn(w http.ResponseWriter, req *http.Request) error {
	log.Println("Sign In")
	paraDef := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rawbody, &paraDef); err != nil {
		return err
	}

	if len(paraDef.Username) < 3 || len(paraDef.Password) < 8 {
		return fmt.Errorf("wrong user or password")
	}

	rspdata := struct {
		Valid bool
		Token string
	}{
		Valid: false,
	}

	return util.WriteJsonResp(w, rspdata)
}
