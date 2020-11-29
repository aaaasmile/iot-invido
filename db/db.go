package db

import (
	"fmt"
	"time"

	"github.com/aaaasmile/iot-invido/web/idl"
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

func (conn *InfluxDbConn) InsertSensorData(sendId string, useDeltaTime bool, sensState *idl.SensorState) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: conn.dbHost,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  conn.dbName,
		Precision: "s",
	})
	if err != nil {
		return err
	}
	fmt.Println("** Batch point is ", bp)

	// // Create a point and add to batch
	// tags := map[string]string{"productView": productMeasurement["ProductName"].(string)}
	// fields := productMeasurement

	pt, err := client.NewPoint("products", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// // Write the batch
	// if err := c.Write(bp); err != nil {
	// 	return err
	// }

	// // Close client resources
	// if err := c.Close(); err != nil {
	// 	return err
	// }

	return nil
}
