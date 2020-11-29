package db

import (
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	client "github.com/influxdata/influxdb1-client/v2"
)

type InfluxConn struct {
	dbName string
	dbHost string
}

func NewInfluxConn(dbname, host string) *InfluxConn {
	con := InfluxConn{
		dbName: dbname,
		dbHost: host,
	}
	return &con
}


// Insert saves points to database
func Insert(productMeasurement map[string]interface{}) error {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: ,
	})
	if err != nil {
		return err
	}
	defer c.Close()

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  MyDB,
		Precision: "s",
	})
	if err != nil {
		return err
	}

	// Create a point and add to batch
	tags := map[string]string{"productView": productMeasurement["ProductName"].(string)}
	fields := productMeasurement

	pt, err := client.NewPoint("products", tags, fields, time.Now())
	if err != nil {
		return err
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		return err
	}

	// Close client resources
	if err := c.Close(); err != nil {
		return err
	}

	return nil
}
