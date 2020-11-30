package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aaaasmile/iot-invido/web/idl"
	"github.com/aaaasmile/iot-invido/web/iot/sensor"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDbConn struct {
	bucketName string
	dbHost     string
	token      string
	org        string
}

func NewInfluxConn(info *idl.Influx) *InfluxDbConn {
	con := InfluxDbConn{
		bucketName: info.BucketName,
		dbHost:     info.DbHost,
		org:        info.Org,
		token:      info.Token,
	}
	return &con
}

func (conn *InfluxDbConn) InsertSensorData(name string, useDeltaTime bool, prevTimestamp time.Time, senSt *sensor.SensorState) error {
	client := influxdb2.NewClient(conn.dbHost, conn.token)
	defer client.Close()

	writeAPI := client.WriteAPI(conn.org, conn.bucketName)

	tags := map[string]string{
		"air-quality-class": senSt.GetAirQualityTag(),
		"place":             senSt.Place,
		"sensor-id":         senSt.SensorID,
	}

	fields := senSt.GetInterfaceMap()
	ts := time.Now()
	if useDeltaTime && senSt.TimeStamp.After(prevTimestamp) {
		ts = prevTimestamp.Add(senSt.TimeStamp.Sub(prevTimestamp))
	}
	senSt.TimeStamp = ts

	p := influxdb2.NewPoint(name,
		tags,
		fields,
		time.Now())
	// write point asynchronously
	writeAPI.WritePoint(p)
	// Flush writes
	writeAPI.Flush()

	log.Println("Batch point inserted ", p)
	return nil
}

func (conn *InfluxDbConn) FetchData(name string) ([]sensor.SensorState, error) {
	client := influxdb2.NewClient(conn.dbHost, conn.token)
	defer client.Close()

	list := []sensor.SensorState{}

	query := `
	from(bucket:"%s")
		|> range(start: -2h)
		|> filter(fn: (r) => r._measurement == "%s")
		|> filter(fn: (r) => r["_field"] == "iaq")
		|> filter(fn: (r) => r["place"] == "Home")
		|> filter(fn: (r) => r["sensor-id"] == "Test")
	`
	query = fmt.Sprintf(query, conn.bucketName, name)
	queryAPI := client.QueryAPI(conn.org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Println("Error on query compile")
		return nil, err
	}
	for result.Next() {
		// Notice when group key has changed
		// if result.TableChanged() {
		// 	fmt.Printf("table: %s\n", result.TableMetadata().String())
		// }
		// Access data
		fmt.Printf("*** value: %v\n", result.Record().Values())
		unk := result.Record().Value()
		if fv, ok := unk.(float64); ok {
			ss := sensor.SensorState{
				Iaq:      float32(fv),
				IaqClass: result.Record().ValueByKey("air-quality-class").(string),
			}
			list = append(list, ss)
		} else {
			fmt.Println("** Not reco", unk)
		}
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("query parsing error: %v", result.Err().Error())
	}
	return list, nil
}
