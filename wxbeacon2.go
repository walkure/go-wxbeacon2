package wxbeacon2

import (
	"encoding/binary"
	"log"

	"github.com/bettercap/gatt"
)

var targetDeviceId string

type WxCallback func(data interface{})

var wxCbFunc WxCallback

type commonData struct {
	DeviceId     string
	Sequence     byte
	Temp         float32
	Humid        float32
	AmbientLight uint16
	UVIndex      float32
	Pressure     float32
	SoundNoise   float32
	VBattery     float32
}

type WxIMData struct {
	commonData
	AccelerationX uint16
	AccelerationY uint16
	AccelerationZ uint16
}

type WxEPData struct {
	commonData
	DisconfortIndex float32
	HeatStroke      float32
}

//cast pattern : https://play.golang.org/p/n1_YO_t2gYK

func parseIM(deviceId string, data []byte) WxIMData {
	parsed := WxIMData{}

	parsed.DeviceId = deviceId
	parsed.Sequence = data[2]
	parsed.Temp = float32(int16(binary.LittleEndian.Uint16(data[3:5]))) / 100
	parsed.Humid = float32(binary.LittleEndian.Uint16(data[5:7])) / 100
	parsed.AmbientLight = binary.LittleEndian.Uint16(data[7:9])
	parsed.UVIndex = float32(binary.LittleEndian.Uint16(data[9:11])) / 100
	parsed.Pressure = float32(binary.LittleEndian.Uint16(data[11:13])) / 10
	parsed.SoundNoise = float32(binary.LittleEndian.Uint16(data[13:15])) / 100

	parsed.AccelerationX = binary.LittleEndian.Uint16(data[15:17])
	parsed.AccelerationY = binary.LittleEndian.Uint16(data[17:19])
	parsed.AccelerationZ = binary.LittleEndian.Uint16(data[19:21])

	parsed.VBattery = float32(int16(data[21])+100) / 100

	return parsed
}

func parseEP(deviceId string, data []byte) WxEPData {
	parsed := WxEPData{}

	parsed.DeviceId = deviceId
	parsed.Sequence = data[2]
	parsed.Temp = float32(int16(binary.LittleEndian.Uint16(data[3:5]))) / 100
	parsed.Humid = float32(binary.LittleEndian.Uint16(data[5:7])) / 100
	parsed.AmbientLight = binary.LittleEndian.Uint16(data[7:9])
	parsed.UVIndex = float32(binary.LittleEndian.Uint16(data[9:11])) / 100
	parsed.Pressure = float32(binary.LittleEndian.Uint16(data[11:13])) / 10
	parsed.SoundNoise = float32(binary.LittleEndian.Uint16(data[13:15])) / 100

	parsed.DisconfortIndex = float32(binary.LittleEndian.Uint16(data[15:17])) / 100
	parsed.HeatStroke = float32(binary.LittleEndian.Uint16(data[17:19])) / 100

	parsed.VBattery = float32(int16(data[21])+100) / 100

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

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {

	if p.ID() != targetDeviceId {
		return
	}

	if wxCbFunc == nil {
		return
	}

	switch p.Name() {
	case "EP":
		wxCbFunc(parseEP(p.ID(), a.ManufacturerData))
	case "IM":
		wxCbFunc(parseIM(p.ID(), a.ManufacturerData))
	default:
		log.Fatalf("Unknown Name:%s", p.Name())
	}
}

func WaitForReceiveData(deviceId string, cb WxCallback) error {

	d, err := gatt.NewDevice()
	if err != nil {
		return err
	}

	targetDeviceId = deviceId
	wxCbFunc = cb

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)

	return nil
}
