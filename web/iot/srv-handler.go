package iot

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/aaaasmile/iot-invido/conf"
	"github.com/aaaasmile/iot-invido/web/idl"
	"github.com/aaaasmile/iot-invido/web/iot/sensor"
)

type PageCtx struct {
	RootUrl    string
	Buildnr    string
	VueLibName string
}

func getURLForRoute(uri string) string {
	arr := strings.Split(uri, "/")
	//fmt.Println("split: ", arr, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		ss := arr[i]
		if ss != "" {
			if !strings.HasPrefix(ss, "?") {
				//fmt.Printf("Url for route is %s\n", ss)
				return ss
			}
		}
	}
	return uri
}

func APiHandler(w http.ResponseWriter, req *http.Request) {
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

func handleGet(w http.ResponseWriter, req *http.Request) error {
	u, _ := url.Parse(req.RequestURI)
	log.Println("GET requested ", u)

	pagectx := PageCtx{
		RootUrl:    conf.Current.RootURLPattern,
		Buildnr:    idl.Buildnr,
		VueLibName: conf.Current.VueLibName,
	}
	templName := "templates/vue/index.html"

	tmplIndex := template.Must(template.New("AppIndex").ParseFiles(templName))

	err := tmplIndex.ExecuteTemplate(w, "base", pagectx)
	if err != nil {
		return err
	}
	return nil
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

func InitFromConfig(debug bool) error {
	// todo open the database
	log.Println("Handler initialized", debug)
	return nil
}

func init() {
	statusCh := make(chan *sensor.SensorState)
	go listenStatus(statusCh)
}
