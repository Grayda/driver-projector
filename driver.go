package main

import (
	"github.com/ninjasphere/go-ninja/api"
	"github.com/ninjasphere/go-ninja/logger"
	"github.com/ninjasphere/go-ninja/model"
)

const (
	driverName = "photography.davidgray.projector"
)

var log = logger.GetLogger(driverName)
var info = ninja.LoadModuleInfo("./package.json")

type WemoDeviceContext struct {
	Info       *wemo.DeviceInfo
	Device     *wemo.Device
	deviceInfo *model.Device
	driver     ninja.Driver
}

func (w *WemoDeviceContext) GetDeviceInfo() *model.Device {
	return w.deviceInfo
}

func (w *WemoDeviceContext) GetDriver() ninja.Driver {
	return w.driver
}

func (w *WemoDeviceContext) SetEventHandler(sendEvent func(event string, payload interface{}) error) {
}

type WemoDriver struct {
	conn      *ninja.Connection
	sendEvent func(event string, payload interface{}) error
}

func NewWemoDriver() (*WemoDriver, error) {
	conn, err := ninja.Connect(driverName)
	if err != nil {
		log.HandleError(err, "Could not connect to MQTT")
		return nil, err
	}

	driver := &WemoDriver{
		conn: conn,
	}

	err = conn.ExportDriver(driver)
	if err != nil {
		log.Fatalf("Failed to export Wemo driver: %s", err)
	}

	return driver, nil
}

func (d *WemoDriver) Start(x interface{}) error {
	log.Infof("Start method on Wemo driver called")

	return d.startDiscovery()
}
