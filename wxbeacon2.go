package wxbeacon2

import (
	"encoding/binary"
	"log"

	"github.com/bettercap/gatt"
)

type WxCallback func(data interface{})


type WxBeaconDevice struct{
	wxCbFunc WxCallback
	device gatt.Device
	targetDeviceId string
}

type commonData struct {
	DeviceId     string
	Sequence     byte
	Temp         float64
	Humid        float64
	AmbientLight uint16
	UVIndex      float64
	Pressure     float64
	SoundNoise   float64
	VBattery     float64
}

type WxIMData struct {
	commonData
	AccelerationX uint16
	AccelerationY uint16
	AccelerationZ uint16
}

type WxEPData struct {
	commonData
	DisconfortIndex float64
	HeatStroke      float64
}

//cast pattern : https://play.golang.org/p/n1_YO_t2gYK

func parseIM(deviceId string, data []byte) WxIMData {
	parsed := WxIMData{}

	parsed.DeviceId = deviceId
	parsed.Sequence = data[2]
	parsed.Temp = float64(int16(binary.LittleEndian.Uint16(data[3:5]))) / 100
	parsed.Humid = float64(binary.LittleEndian.Uint16(data[5:7])) / 100
	parsed.AmbientLight = binary.LittleEndian.Uint16(data[7:9])
	parsed.UVIndex = float64(binary.LittleEndian.Uint16(data[9:11])) / 100
	parsed.Pressure = float64(binary.LittleEndian.Uint16(data[11:13])) / 10
	parsed.SoundNoise = float64(binary.LittleEndian.Uint16(data[13:15])) / 100

	parsed.AccelerationX = binary.LittleEndian.Uint16(data[15:17])
	parsed.AccelerationY = binary.LittleEndian.Uint16(data[17:19])
	parsed.AccelerationZ = binary.LittleEndian.Uint16(data[19:21])

	parsed.VBattery = float64(int16(data[21])+100) / 100

	return parsed
}

func parseEP(deviceId string, data []byte) WxEPData {
	parsed := WxEPData{}

	parsed.DeviceId = deviceId
	parsed.Sequence = data[2]
	parsed.Temp = float64(int16(binary.LittleEndian.Uint16(data[3:5]))) / 100
	parsed.Humid = float64(binary.LittleEndian.Uint16(data[5:7])) / 100
	parsed.AmbientLight = binary.LittleEndian.Uint16(data[7:9])
	parsed.UVIndex = float64(binary.LittleEndian.Uint16(data[9:11])) / 100
	parsed.Pressure = float64(binary.LittleEndian.Uint16(data[11:13])) / 10
	parsed.SoundNoise = float64(binary.LittleEndian.Uint16(data[13:15])) / 100

	parsed.DisconfortIndex = float64(binary.LittleEndian.Uint16(data[15:17])) / 100
	parsed.HeatStroke = float64(binary.LittleEndian.Uint16(data[17:19])) / 100

	parsed.VBattery = float64(int16(data[21])+100) / 100

	return parsed
}

func onStateChanged(d gatt.Device, s gatt.State) {
	switch s {
	case gatt.StatePoweredOn:
		// allow duplicate
		d.Scan([]gatt.UUID{}, true)
		return
	default:
		d.StopScanning()
	}
}

func (d WxBeaconDevice)onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {

	if p.ID() != d.targetDeviceId {
		return
	}

	if d.wxCbFunc == nil {
		return
	}

	switch p.Name() {
	case "EP":
		d.wxCbFunc(parseEP(p.ID(), a.ManufacturerData))
	case "IM":
		d.wxCbFunc(parseIM(p.ID(), a.ManufacturerData))
	default:
		log.Fatalf("Unknown Name:%s", p.Name())
	}
}

func WaitForReceiveData(deviceId string, cb WxCallback) (*WxBeaconDevice,error) {

	d, err := gatt.NewDevice()
	if err != nil {
		return nil,err
	}
	dev := WxBeaconDevice{}

	dev.targetDeviceId = deviceId
	dev.wxCbFunc = cb
	dev.device = d

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(dev.onPeriphDiscovered))
	d.Init(onStateChanged)

	return &dev,nil
}

func (d WxBeaconDevice) Stop() error {
	return d.device.Stop()
}
