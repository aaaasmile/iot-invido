#!/bin/bash

echo "Builds app"
go build -o iot-invido.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target iotinvido -outdir ~/app/go/iot-invido/zips/
cd ~/app/go/iot-invido/

echo "update the service"
./update-service.sh

echo "Ready to fly"