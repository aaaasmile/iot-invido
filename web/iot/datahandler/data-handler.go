package datahandler

import (
	"log"
	"net/http"
	"time"

	"github.com/aaaasmile/iot-invido/db"
	"github.com/aaaasmile/iot-invido/util"
	"github.com/aaaasmile/iot-invido/web/idl"
	"github.com/aaaasmile/iot-invido/web/iot/sensor"
)

type HandleData struct {
	Influx *idl.Influx
}

type RespData struct {
	Status   string               `json:"status"`
	DataView []sensor.SensorState `json:"dataview"`
}

func (hd *HandleData) HandleTestInsertLine(w http.ResponseWriter, req *http.Request) error {
	log.Println("Insert test data")
	conn := db.NewInfluxConn(hd.Influx)
	sensState := sensor.SensorState{SensorID: "Test", Place: "Home"}
	sensState.SetRandomData()

	prevTs := time.Now()
	if err := conn.InsertSensorData("SimBM680", false, prevTs, &sensState); err != nil {
		return err
	}

	list := []sensor.SensorState{sensState}
	rspdata := RespData{
		Status:   "OK",
		DataView: list,
	}

	return util.WriteJsonResp(w, rspdata)
}

func (hd *HandleData) HandleFetchData(w http.ResponseWriter, req *http.Request) error {
	log.Println("Fetch data")
	conn := db.NewInfluxConn(hd.Influx)

	list, err := conn.FetchData("SimBM680")
	if err != nil {
		return err
	}

	rspdata := RespData{
		Status:   "OK",
		DataView: list,
	}
	return util.WriteJsonResp(w, rspdata)
}
