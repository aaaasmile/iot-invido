package db

import (
	"fmt"
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

	fmt.Println("** Batch point inserted ", p)

	return nil
}
