package datahandler

import (
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

func (hd *HandleData) HandleTestInsertLine(w http.ResponseWriter, req *http.Request) error {
	conn := db.NewInfluxConn(hd.Influx.DbHost, hd.Influx.DbName)
	sensState := sensor.SensorState{SensorID: "Test", Place: "Home"}
	sensState.SetRandomData()

	prevTs := time.Now()
	if err := conn.InsertSensorData("SimBM680", false, prevTs, &sensState); err != nil {
		return err
	}

	list := []sensor.SensorState{sensState}
	rspdata := struct {
		Status   string               `json:"status"`
		DataView []sensor.SensorState `json:"dataview"`
	}{
		Status:   "OK",
		DataView: list,
	}

	return util.WriteJsonResp(w, rspdata)
}
