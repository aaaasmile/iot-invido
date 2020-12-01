package datahandler

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func (hd *HandleData) HandlePubData(w http.ResponseWriter, req *http.Request) error {
	rawbody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	lines := string(rawbody)
	log.Println("Body is: ", lines)
	ll := strings.Split(lines, "\n")
	for _, line := range ll {
		sensState := sensor.SensorState{SensorID: "BMI680-1", Place: "Home"}
		fieldsArr := strings.Split(line, ",")
		for _, fieldKV := range fieldsArr {
			pair := strings.Split(fieldKV, ":")
			if len(pair) == 2 {
				myVal := strings.Trim(pair[1], " ")
				myKey := strings.Trim(pair[0], " ")
				switch myKey {
				case "TS":
					sensState.TimeStamp = time.Now()
				case "TEMP-RAW":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.TempRaw = float32(value)
				case "PRES":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Press = float32(value)
				case "HUMI-RAW":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.HumiRaw = float32(value)
				case "GASO":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Gaso = float32(value)
				case "IAQ":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Iaq = float32(value)
				case "IAQA":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Iaqacc = float32(value)
				case "TEMP":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Temp = float32(value)
				case "HUMY":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Humidity = float32(value)
				case "CO2":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Co2 = float32(value)
				case "VOC":
					value, _ := strconv.ParseFloat(myVal, 32)
					sensState.Voc = float32(value)
				}
			}
		}
		log.Println("Recognized sensor data: ", sensState)
		conn := db.NewInfluxConn(hd.Influx)
		if err := conn.InsertSensorData("SimBM680", false, sensState.TimeStamp, &sensState); err != nil {
			return err
		}

	}
	// Body is:  TS: 140255, TEMP-RAW: 21.41, PRES: 100174.00, HUMI-RAW: 47.32, GASO: 84509.00, IAQ: 25.00, IAQA: 0, TEMP: 21.35, HUMY: 47.46, CO2: 500.00, VOC: 0.50
	return nil
}
