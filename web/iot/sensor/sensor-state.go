package sensor

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

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
	IaqClass  string    `json:"iaqclass"`
}

func (ss *SensorState) CalculateAirQualityTag() string {
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
	ss.IaqClass = iaqstr
	return iaqstr
}

func (ss *SensorState) GetFieldsMap() map[string]interface{} {
	fields := map[string]interface{}{
		"temperature":  ss.Temp,
		"pressure":     ss.Press,
		"humidity":     ss.Humidity,
		"gasraw":       ss.Gaso,
		"iaq":          ss.Iaq,
		"iaqaccurancy": ss.Iaqacc,
		"co2":          ss.Co2,
		"voc":          ss.Voc,
		"iaqclass":     ss.IaqClass,
	}
	return fields
}

func (ss *SensorState) GetTagsMap() map[string]string {
	tags := map[string]string{
		"place":     ss.Place,
		"sensor-id": ss.SensorID,
	}
	return tags
}

func (ss *SensorState) SetRandomData() {
	ss.Temp = 18.2 + rand.Float32()*10.0
	ss.Press = 1024 + rand.Float32()*10.0
	ss.Humidity = 43.6 + rand.Float32()*30.0
	ss.Gaso = 34009.8
	ss.Iaq = 10 + rand.Float32()*300.0
	ss.Iaqacc = 3
	ss.Co2 = 500 + rand.Float32()*400.0
	ss.Voc = 1.34 + +rand.Float32()*13.0
	ss.CalculateAirQualityTag()
}

func (ss *SensorState) SetMembersFromDBMapValues(mapValues map[string]interface{}) error {
	// tt, err := time.Parse(time.RFC3339, mapValues["_time"])
	// if err != nil {
	// 	log.Println("Error on parsing timestamp")
	// 	return err
	// }
	var tt time.Time
	var ok bool
	var fval float32
	var strval string

	ttUnk := mapValues["_time"]
	if tt, ok = ttUnk.(time.Time); !ok {
		log.Println("Error on parsing timestamp")
		return fmt.Errorf("Error on recognize timestamp. %v", ttUnk)
	}
	ss.TimeStamp = tt

	field := mapValues["_field"]
	valueUnk := mapValues["_value"]
	placestr, err := getValueAsString(mapValues["place"])
	if err != nil {
		return fmt.Errorf("Expect string in tag place. Error %v", err)
	}
	ss.Place = placestr

	sensorstr, err := getValueAsString(mapValues["sensor-id"])
	if err != nil {
		return fmt.Errorf("Expect string in tag sensor-id. Error %v", err)
	}
	ss.SensorID = sensorstr

	if fv, ok := valueUnk.(float64); ok {
		fval = float32(fv)
	} else if str, ok := valueUnk.(string); ok {
		strval = str
	} else {
		return fmt.Errorf("Value not recognized: %v", valueUnk)
	}

	switch field.(string) {
	case "temperature":
		ss.Temp = fval
		ss.TempRaw = fval
	case "pressure":
		ss.Press = fval
	case "humidity":
		ss.HumiRaw = fval
		ss.Humidity = fval
	case "gasraw":
		ss.Gaso = fval
	case "iaq":
		ss.Iaq = fval
	case "iaqaccurancy":
		ss.Iaqacc = fval
	case "co2":
		ss.Co2 = fval
	case "voc":
		ss.Voc = fval
	case "iaqclass":
		ss.IaqClass = strval
	}
	return nil
}

func getValueAsString(val interface{}) (string, error) {
	if v, ok := val.(string); ok {
		return v, nil
	}
	return "", fmt.Errorf("The value interface is not a string")
}
