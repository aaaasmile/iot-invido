package idl

import "golang.org/x/text/message"

var (
	Appname = "iot-invido"
	Buildnr = "00.03.01.20220203-00"
	Printer *message.Printer
)

type Influx struct {
	BucketName string
	DbHost     string
	Org        string
	Token      string
}
