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
	tags := senSt.GetTagsMap()

	fields := senSt.GetFieldsMap()
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

	log.Println("Batch point inserted ", p.Name())
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
		|> filter(fn: (r) => r["_field"] == "temperature" or r["_field"] == "iaq")
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
		//fmt.Printf("*** value: %v\n", result.Record().Values())
		ss := sensor.SensorState{}
		mapValues := result.Record().Values()
		//la mapValues è una mappa di coppie valori key/value
		// Per un dato punto (chiave _time) si hanno unsa serie di coppie dove si ha il valore del campo, il suo nome e i vari tags.
		// Per esempio un'istanza di mapValues è:
		//    map[_field:temperature _measurement:SimBM680 _start:2020-12-01 19:08:12.5700177 +0000 UTC _stop:2020-12-01 21:08:12.5700177 +0000 UTC _time:2020-12-01
		//    20:36:22.501973 +0000 UTC _value:23.70146942138672 place:Home result:_result sensor-id:Test table:1]
		// Quindi map["_field"] ti dice qual è il campo del mapValues in questione
		//        map["_time"] ti dice il timestamp
		//        map["_value"] ti dice il valore
		// tutte le altre chiavi ti danno delle info in più.

		// for k, v := range mapValues {
		// 	fmt.Println("*** k: ", k)
		// 	fmt.Println("*** v: ", v)
		// }
		//fmt.Println("*** values: ", mapValues)
		ss.SetMembersFromDBMapValues(mapValues)
		fmt.Println("*** ss is : ", ss)
		// if fv, ok := mapValues["iaq"]; ok {
		// 	ss.Iaq = float32(fv.(float64))
		// }
		// if fv, ok := mapValues["iaqclass"]; ok {
		// 	ss.IaqClass = fv.(string)
		// }
		// if len(list) < 10 {
		list = append(list, ss)
		// }
	}

	if result.Err() != nil {
		return nil, fmt.Errorf("query parsing error: %v", result.Err().Error())
	}
	return list, nil
}
