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
	writeAPI.WritePoint(p)
	writeAPI.Flush()

	log.Println("Batch asynch point inserted ", p.Name())
	return nil
}

func (conn *InfluxDbConn) FetchData(name string) ([]*sensor.SensorState, error) {
	client := influxdb2.NewClient(conn.dbHost, conn.token)
	defer client.Close()

	// Sui campi da ritornare, meno sono e meglio è. La mia UI ne richiede diversi.
	// Questa è una query che ritorna gli ultimi 20 punti da quando la si chiama
	query := `
	from(bucket:"%s")
		|> range(start: -2h)
		|> filter(fn: (r) => r._measurement == "%s")
		|> filter(fn: (r) => r["_field"] == "temperature" or r["_field"] == "iaq" 
				 or r["_field"] == "pressure" or r["_field"] == "humidity" 
				 or r["_field"] == "co2" or r["_field"] == "iaqaccurancy")
		|> sort(columns:["_time"], desc: true)
		|> limit(n: 20)
	`
	query = fmt.Sprintf(query, conn.bucketName, name)
	queryAPI := client.QueryAPI(conn.org)
	result, err := queryAPI.Query(context.Background(), query)
	if err != nil {
		log.Println("Error on query compile")
		return nil, err
	}

	list := []*sensor.SensorState{}
	ix := 0
	currField := ""
	insert := true
	var ss *sensor.SensorState
	for result.Next() {

		mapValues := result.Record().Values()
		//la mapValues è una mappa di coppie valori key/value
		// Per un dato punto (chiave _time) si hanno una serie di coppie dove si ha il valore del campo, il suo nome e i vari tags.
		// Per esempio un'istanza di mapValues per undato field ad un dato timestamp è:
		//    map[_field:temperature _measurement:SimBM680 _start:2020-12-01 19:08:12.5700177 +0000 UTC _stop:2020-12-01 21:08:12.5700177 +0000 UTC _time:2020-12-01
		//    20:36:22.501973 +0000 UTC _value:23.70146942138672 place:Home result:_result sensor-id:Test table:1]
		// Quindi map["_field"] ti dice qual è il campo del mapValues in questione
		//        map["_time"] ti dice il timestamp
		//        map["_value"] ti dice il valore
		// tutte le altre chiavi ti danno delle info in più.
		// Esempio: Se ho 10 timestamp nel mio range, ed ho richiesto 3 campi,
		//          questo Next() verrà chiamato 3 x 10 ogni volta con un campo diverso, ma usando sempre gli stessi timestamps
		//          È un punto che ho fatto fatica a capire, in quanto salvando un'instanza di SensorState, improvvisamente mi ritrovo con molti più record.
		// Risulta evidente che mostrare più campi per un timestamp non è ideale e risulta in una quantità maggiore di dati da combinare.

		// for k, v := range mapValues {
		// 	fmt.Println("*** k: ", k)
		// 	fmt.Println("*** v: ", v)
		// }
		//fmt.Println("*** values: ", mapValues)
		if str, ok := mapValues["_field"].(string); ok {
			if currField == "" {
				currField = str
				insert = true
				ix = 0
			} else if currField != str {
				insert = false
				ix = 0
				currField = str
				log.Println("Processing the field", str)
			}
		} else {
			return nil, fmt.Errorf("Unable to read the field caption")
		}

		if insert {
			ss = &sensor.SensorState{}
		} else {
			ss = list[ix]
		}

		ss.SetMembersFromDBMapValues(mapValues)
		//fmt.Println("*** ss is : ", ss)
		if insert {
			list = append(list, ss)
		} else {
			ix++
		}
	}
	log.Println("Collected points : ", len(list))

	if result.Err() != nil {
		return nil, fmt.Errorf("query parsing error: %v", result.Err().Error())
	}
	return list, nil
}
