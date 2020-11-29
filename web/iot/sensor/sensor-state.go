package sensor

import "time"

type SensorState struct {
	SensorID  string    `json:"sensorid"`
	Place     string    `json:"place"`
	TempRaw   float32   `json:"tempraw"`
	Press     float32   `json:"press"`
	HumiRaw   float32   `json:"humiraw"`
	Gaso      float32   `json:"gaso"`
	Iaq       float32   `json:"iaq"`
	Iaqacc    float32   `json:"iaqacc"`
	Temp      float32   `json:"temp"`
	Humidity  float32   `json:"humy"`
	Co2       float32   `json:"co2"`
	Voc       float32   `json:"voc"`
	TimeStamp time.Time `json:"timeStamp"`
}

func (ss *SensorState) GetAirQualityTag() string {
	iaqstr := ""
	iaq := ss.Iaq
	if iaq > 301 {
		iaqstr = "Hazardous"
	} else if iaq > 250 && iaq <= 300 {
		iaqstr = "Very Unhealthy"
	} else if iaq > 200 && iaq <= 250 {
		iaqstr = "More than Unhealthy"
	} else if iaq > 150 && iaq <= 200 {
		iaqstr = "Unhealthy"
	} else if iaq > 100 && iaq <= 150 {
		iaqstr = "Unhealthy for Sensitive Groups"
	} else if iaq > 50 && iaq <= 100 {
		iaqstr = "Moderate"
	} else if iaq >= 00 && iaq <= 50 {
		iaqstr = "Good"
	}
	return iaqstr
}

func (ss *SensorState) GetInterfaceMap() map[string]interface{} {
	fields := map[string]interface{}{
		"sensorid":     ss.SensorID,
		"place":        ss.Place,
		"temperature":  ss.Temp,
		"pressure":     ss.Press,
		"humidity":     ss.Humidity,
		"gasraw":       ss.Gaso,
		"iaq":          ss.Iaq,
		"iaqaccurancy": ss.Iaqacc,
		"co2":          ss.Co2,
		"voc":          ss.Voc,
	}
	return fields
}

func (ss *SensorState) SetRandomData() {
	ss.Temp = 25.7
	ss.Press = 1024
	ss.Humidity = 43.6
	ss.Gaso = 34009.8
	ss.Iaq = 71.6
	ss.Iaqacc = 3
	ss.Co2 = 500
	ss.Voc = 1.34
}
