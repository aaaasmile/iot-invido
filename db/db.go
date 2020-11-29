package db

import (
	"fmt"
	"time"

	"github.com/aaaasmile/iot-invido/web/iot/sensor"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxDbConn struct {
	dbName string
	dbHost string
}

func NewInfluxConn(host, dbname string) *InfluxDbConn {
	con := InfluxDbConn{
		dbName: dbname,
		dbHost: host,
	}
	return &con
}

func (conn *InfluxDbConn) InsertSensorData(name string, useDeltaTime bool, prevTimestamp time.Time, senSt *sensor.SensorState) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conn.dbHost,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.dbName,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	fmt.Println("** Batch point is ", bp)

	tags := map[string]string{"AirQuality": senSt.GetAirQualityTag()}
	fields := senSt.GetInterfaceMap()
	ts := time.Now()
	if useDeltaTime && senSt.TimeStamp.After(prevTimestamp) {
		ts = prevTimestamp.Add(senSt.TimeStamp.Sub(prevTimestamp))
	}
	pt, err := client.NewPoint(name, tags, fields, ts)
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	if err := c.Write(bp); err != nil {
		return err
	}

	return nil
}
