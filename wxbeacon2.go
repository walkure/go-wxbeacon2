package wxbeacon2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/walkure/gatt"
)

type WxCallback func(data interface{})

type Device struct {
	wxCbFunc       WxCallback
	device         gatt.Device
	targetDeviceId string
}

type commonData struct {
	RSSI         int
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

func (d WxIMData) String() string {
	sb := strings.Builder{}
	sb.WriteString("Type:IM DeviceId:")
	sb.WriteString(d.DeviceId)
	sb.WriteString(fmt.Sprintf(" RSSI:%d", d.RSSI))
	sb.WriteString(fmt.Sprintf(" Sequence:%d", d.Sequence))
	sb.WriteString(fmt.Sprintf(" Temp:%g", d.Temp))
	sb.WriteString(fmt.Sprintf(" Humid:%g", d.Humid))
	sb.WriteString(fmt.Sprintf(" AmbientLight:%d", d.AmbientLight))
	sb.WriteString(fmt.Sprintf(" UVIndex:%g", d.UVIndex))
	sb.WriteString(fmt.Sprintf(" Pressure:%g", d.Pressure))
	sb.WriteString(fmt.Sprintf(" SoundNoise:%g", d.SoundNoise))
	sb.WriteString(fmt.Sprintf(" AccelerationX:%d", d.AccelerationX))
	sb.WriteString(fmt.Sprintf(" AccelerationY:%d", d.AccelerationY))
	sb.WriteString(fmt.Sprintf(" AccelerationZ:%d", d.AccelerationZ))
	sb.WriteString(fmt.Sprintf(" VBattery:%g", d.VBattery))
	return sb.String()
}

type WxEPData struct {
	commonData
	DisconfortIndex float64
	HeatStroke      float64
}

func (d WxEPData) String() string {
	sb := strings.Builder{}
	sb.WriteString("Type:EP DeviceId:")
	sb.WriteString(d.DeviceId)
	sb.WriteString(fmt.Sprintf(" RSSI:%d", d.RSSI))
	sb.WriteString(fmt.Sprintf(" Sequence:%d", d.Sequence))
	sb.WriteString(fmt.Sprintf(" Temp:%g", d.Temp))
	sb.WriteString(fmt.Sprintf(" Humid:%g", d.Humid))
	sb.WriteString(fmt.Sprintf(" AmbientLight:%d", d.AmbientLight))
	sb.WriteString(fmt.Sprintf(" UVIndex:%g", d.UVIndex))
	sb.WriteString(fmt.Sprintf(" Pressure:%g", d.Pressure))
	sb.WriteString(fmt.Sprintf(" SoundNoise:%g", d.SoundNoise))
	sb.WriteString(fmt.Sprintf(" DisconfortIndex:%g", d.DisconfortIndex))
	sb.WriteString(fmt.Sprintf(" HeatStroke:%g", d.HeatStroke))
	sb.WriteString(fmt.Sprintf(" VBattery:%g", d.VBattery))
	return sb.String()
}

//cast pattern : https://play.golang.org/p/n1_YO_t2gYK

func parseIM(deviceId string, rssi int, data []byte) WxIMData {
	parsed := WxIMData{}

	parsed.RSSI = rssi
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

func parseEP(deviceId string, rssi int, data []byte) WxEPData {
	parsed := WxEPData{}

	parsed.RSSI = rssi
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

func (dev Device) onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {

	if strings.ToUpper(p.ID()) != dev.targetDeviceId {
		//log.Printf("ignore device:%q %q", p.ID(),p.Name())
		return
	}

	if dev.wxCbFunc == nil {
		return
	}

	switch p.Name() {
	case "EP":
		dev.wxCbFunc(parseEP(p.ID(), rssi, a.ManufacturerData))
	case "IM":
		dev.wxCbFunc(parseIM(p.ID(), rssi, a.ManufacturerData))
	default:
		log.Fatalf("Unknown Name:%s", p.Name())
	}
}

func NewDevice(deviceId string, cb WxCallback) *Device {
	dev := Device{}

	dev.targetDeviceId = strings.ToUpper(deviceId)
	dev.wxCbFunc = cb
	return &dev
}

func (dev *Device) WaitForReceiveData() error {

	if dev == nil || dev.targetDeviceId == "" || dev.wxCbFunc == nil {
		return nil
	}
	d, err := gatt.NewDevice(gatt.LnxSetScanMode(false))
	if err != nil {
		return err
	}

	dev.device = d

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(dev.onPeriphDiscovered))
	d.Init(onStateChanged)

	return nil
}

func (dev *Device) Stop() error {
	if dev == nil {
		return nil
	}
	if dev.device == nil {
		return errors.New("device not initialized")
	}
	err := dev.device.Stop()
	dev.device = nil
	return err
}
