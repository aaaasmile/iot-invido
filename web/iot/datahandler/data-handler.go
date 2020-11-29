package datahandler

import (
	"net/http"

	"github.com/aaaasmile/iot-invido/db"
	"github.com/aaaasmile/iot-invido/util"
	"github.com/aaaasmile/iot-invido/web/idl"
)

type HandleData struct {
	Influx *idl.Influx
}

func (hd *HandleData) HandleTestInsertLine(w http.ResponseWriter, req *http.Request) error {
	conn := db.NewInfluxConn(hd.Influx.DbHost, hd.Influx.DbName)
	sensState := idl.SensorState{}
	if err := conn.InsertSensorData("SimBM680", false, &sensState); err != nil {
		return err
	}

	list := []idl.SensorState{sensState}
	rspdata := struct {
		Status   string            `json:"status"`
		DataView []idl.SensorState `json:"dataview"`
	}{
		Status:   "OK",
		DataView: list,
	}

	return util.WriteJsonResp(w, rspdata)
}
