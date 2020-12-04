package idl

var (
	Appname = "iot-invido"
	Buildnr = "00.02.02.20201203-00"
)

type Influx struct {
	BucketName string
	DbHost     string
	Org        string
	Token      string
}
