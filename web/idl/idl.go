package idl

import "time"

var (
	Appname = "iot-invido"
	Buildnr = "00.01.02.20201126-00"
)

type SensorState struct {
	TimeStamp time.Time `json:"timeStamp"`
	tempraw   float32   `json:"tempraw"`
	Press     float32   `json:"press"`
	HumiRaw   float32   `json:"humiraw"`
	Gaso      float32   `json:"gaso"`
	Iaq       float32   `json:"iaq"`
	Iaqacc    float32   `json:"iaqacc"`
	Temp      float32   `json:"temp"`
	Humidity  float32   `json:"humy"`
	Co2       float32   `json:"co2"`
	Voc       float32   `json:"voc"`
}

type Influx struct {
	DbName string
	DbHost string
}
