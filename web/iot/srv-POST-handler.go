package iot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/util"
	"github.com/aaaasmile/iot-invido/web/iot/datahandler"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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
	valid, err := validateAPIHeader(req)
	if err != nil {
		return err
	}
	rspdata := struct {
		Valid bool
	}{
		Valid: valid,
	}

	return util.WriteJsonResp(w, rspdata)
}

func validateAPIHeader(req *http.Request) (bool, error) {
	tk := req.Header.Get("x-api-sessiontoken")
	if tk == "" {
		log.Println("Header x-api-sessiontoken is empty")
		return false, nil
	}
	token, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in HMAC")
		}
		return []byte(conf.Current.Influx.Token), nil
	})

	if err != nil {
		return false, err
	}
	if token.Valid {
		cl := token.Claims
		fmt.Println("** Claims is ", cl)
	}

	return true, nil
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

	valid, err := liteDB.CheckUsernamePassword(paraDef.Username, paraDef.Password)
	if err != nil {
		log.Println("Error on password check: ", err)
		return fmt.Errorf("User or password wrong")
	}

	scope := "apicall"
	usergid := uuid.New().String()
	token, err := GenerateJWT(usergid, scope)
	if err != nil {
		return err
	}

	sessMgr.StoreUser(usergid, paraDef.Username, scope)

	rspdata := struct {
		Valid bool
		Token string
	}{
		Valid: valid,
		Token: token,
	}

	return util.WriteJsonResp(w, rspdata)
}

func GenerateJWT(usergid, scope string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["scope"] = scope
	claims["gid"] = usergid
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString([]byte(conf.Current.Influx.Token))

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}
