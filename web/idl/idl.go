package idl

var (
	Appname = "iot-invido"
	Buildnr = "00.01.02.20201126-00"
)

type Influx struct {
	BucketName string
	DbHost     string
	Org        string
	Token      string
}
