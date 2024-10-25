package wxbeacon2

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"strings"

	"github.com/walkure/gatt"
	"github.com/walkure/gatt/logger"
)

type commonData struct {
	// RSSI is the signal strength
	RSSI int
	// DeviceId is the MAC address of the beacon
	DeviceId string
	// Sequence is the sequence number
	Sequence byte
	// Temp is the temperature celcius in 0.01 unit
	Temp float64
	// Humid is the relative humidity in 0.01%
	Humid float64
	// AmbientLight is the ambient light in lx
	AmbientLight uint16
	// UVIndex is the UV index count in 0.01 unit
	UVIndex float64
	// Pressure is the atmospheric pressure in 0.1hPa
	Pressure float64
	// SoundNoise is the sound noise in 0.01dB
	SoundNoise float64
	// VBattery is the battery voltage in 0.01V
	VBattery float64
}

// WxIMData is the data structure of the IM type (General/Limited Broadcaster 1)
type WxIMData struct {
	commonData
	// AccelerationX is the acceleration in X axis
	AccelerationX uint16
	// AccelerationY is the acceleration in Y axis
	AccelerationY uint16
	// AccelerationZ is the acceleration in Z axis
	AccelerationZ uint16
}

// String returns a string representation of the WxIMData
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

// LogValue returns slog.GroupValue of the WxIMData
func (d WxIMData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", "IM"),
		slog.String("deviceId", d.DeviceId),
		slog.Int("rssi", d.RSSI),
		slog.Uint64("sequence", uint64(d.Sequence)),
		slog.Float64("temperature", d.Temp),
		slog.Float64("relativeHumidity", d.Humid),
		slog.Uint64("ambientLight", uint64(d.AmbientLight)),
		slog.Float64("uvIndex", d.UVIndex),
		slog.Float64("pressure", d.Pressure),
		slog.Float64("soundNoise", d.SoundNoise),
		slog.Uint64("accelerationX", uint64(d.AccelerationX)),
		slog.Uint64("accelerationY", uint64(d.AccelerationY)),
		slog.Uint64("accelerationZ", uint64(d.AccelerationZ)),
		slog.Float64("vBattery", d.VBattery),
	)
}

// WxDataSequence returns the sequence number of the WxIMData
func (d WxIMData) WxDataSequence() uint8 {
	return d.Sequence
}

// WxEPData is the data structure of the EP type (General/Limited Broadcaster 2)
type WxEPData struct {
	commonData
	// DisconfortIndex is the disconfort index in 0.01 unit
	DisconfortIndex float64
	// HeatStroke is the heat stroke in 0.01 degree celcius
	HeatStroke float64
}

// String returns a string representation of the WxEPData
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

// LogValue returns slog.GroupValue of the WxEPData
func (d WxEPData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("type", "EP"),
		slog.String("deviceId", d.DeviceId),
		slog.Int("rssi", d.RSSI),
		slog.Uint64("sequence", uint64(d.Sequence)),
		slog.Float64("temperature", d.Temp),
		slog.Float64("relativeHumidity", d.Humid),
		slog.Uint64("ambientLight", uint64(d.AmbientLight)),
		slog.Float64("uvIndex", d.UVIndex),
		slog.Float64("pressure", d.Pressure),
		slog.Float64("soundNoise", d.SoundNoise),
		slog.Float64("disconfortIndex", d.DisconfortIndex),
		slog.Float64("heatStroke", d.HeatStroke),
		slog.Float64("vBattery", d.VBattery),
	)
}

// WxDataSequence returns the sequence number of the WxEPData
func (d WxEPData) WxDataSequence() uint8 {
	return d.Sequence
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

// WxData is the interface of WxIMData and WxEPData
type WxData interface {
	// String is `fmt.Stringer` interface
	String() string
	// LogValue is `slog.LogValuer` interface
	LogValue() slog.Value
	// WxDataSequence returns the sequence number of the WxData
	WxDataSequence() byte
}

var _ WxData = WxIMData{}
var _ WxData = WxEPData{}

// The Company Identifier of OMRON Corporation
const companyID = 0x02d5

// HandleWxBeacon2 returns a callback function for `gatt.PeripheralDiscovered` that can be used to handle the WxBeacon2 device.
// This function is used to parse the `ADV_IND` packet of the WxBeacon2 device.
// Passive scanning is recommended. ( `SCAN_RSP` packet is not supported.)
//
// The deviceId is the device id of the target WxBeacon2 device. if it is empty, all WxBeacon2 devices will be handled.
// The cb function will be called with new goroutine when receives the advertisement packet from a WxBeacon2 device.
// The next function will be called if the device is not a target WxBeacon2 device and can be null.
func HandleWxBeacon2(deviceId string, cb func(d WxData),
	next func(gatt.Peripheral, *gatt.Advertisement, int)) func(gatt.Peripheral, *gatt.Advertisement, int) {
	if cb == nil {
		panic("cb is nil")
	}
	if next == nil {
		next = func(gatt.Peripheral, *gatt.Advertisement, int) {}
	}

	deviceId = strings.ToUpper(deviceId)

	return func(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
		if a.CompanyID != companyID {
			// Another Manufacturer
			next(p, a, rssi)
			return
		}

		devId := strings.ToUpper(p.ID())
		if deviceId != "" && devId != deviceId {
			// Not the target device
			next(p, a, rssi)
			return
		}

		if len(a.ManufacturerData) < 22 {
			// Truncated data
			next(p, a, rssi)
			return
		}

		switch p.Name() {
		case "EP":
			cb(parseEP(devId, rssi, a.ManufacturerData))
			return
		case "IM":
			cb(parseIM(devId, rssi, a.ManufacturerData))
			return
		}

		// Unknown device Name
		next(p, a, rssi)
	}
}

// SetLogger sets the logger for the package
func SetLogger(newLogger *slog.Logger) {
	logger.SetLogger(newLogger)
}
