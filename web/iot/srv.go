package iot

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/iot-invido/db/sqlite"
	"github.com/aaaasmile/iot-invido/web/idl"
	"github.com/aaaasmile/iot-invido/web/iot/sensor"
)

type PageCtx struct {
	RootUrl    string
	Buildnr    string
	VueLibName string
}

var (
	sessMgr = &SessionManager{
		cookieName:  "IotDataCookie",
		sessionsWS:  make(map[string]*SessionCtx),
		maxlifetime: 3600 * 24, // Max session life in seconds
		gidPerUser:  make(map[string]GidInSession),
	}
	funcMap = template.FuncMap{
		"trans": TranslateString,
	}
	liteDB *sqlite.LiteDB
)

func HandleIndex(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	log.Println("Request: ", req.RequestURI)
	var err error
	switch req.Method {
	case "GET":
		err = handleGet(w, req)
	case "POST":
		log.Println("POST on ", req.RequestURI)
		err = handlePost(w, req)
	}
	if err != nil {
		log.Println("Error exec: ", err)
		http.Error(w, fmt.Sprintf("Internal error on execute: %v", err), http.StatusInternalServerError)
	}

	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Service %s total call duration: %v\n", idl.Appname, elapsed)
}

func TranslateString(s string) string {
	return idl.Printer.Sprintf(s)
}

func writeResponse(w http.ResponseWriter, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	wsClients.Broadcast(string(blobresp))
	w.Write(blobresp)
	return nil
}

func writeErrorResponse(w http.ResponseWriter, errorcode int, resp interface{}) error {
	blobresp, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	http.Error(w, string(blobresp), errorcode)
	return nil
}

func listenStatus(statusCh chan *sensor.SensorState) {
	log.Println("Waiting for status")
	for {
		st := <-statusCh
		resp := struct {
			IAQ  string `json:"iaq"`
			Temp string `json:"temp"`
			Humi string `json:"humi"`
			Gas  string `json:"gas"`
			Co2  string `json:"co2"`
		}{}
		log.Println("Status update received ", st)
		blobresp, err := json.Marshal(resp)
		if err != nil {
			log.Println("Error in state relay: ", err)
		} else {
			wsClients.Broadcast(string(blobresp))
		}
	}
}

func InitFromConfig(debug bool, dbsqlitePath string) error {
	log.Println("InitFromConfig. Path, Debug: ", dbsqlitePath, debug)
	liteDB = &sqlite.LiteDB{
		DebugSQL:     debug,
		SqliteDBPath: dbsqlitePath,
	}
	if err := liteDB.OpenSqliteDatabase(); err != nil {
		return err
	}
	InitSession()
	InitWS()
	statusCh := make(chan *sensor.SensorState)
	go listenStatus(statusCh)
	return nil
}
