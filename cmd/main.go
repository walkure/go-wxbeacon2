package main

import (
	"github.com/walkure/go-wxbeacon2"
	"fmt"
	"log"
)

func main() {
	dev := wxbeacon2.NewDevice("ZZ:ZZ:ZZ:ZZ:ZZ:ZZ", printData)
	err := dev.WaitForReceiveData()
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	select {}
}

func printData(data interface{}) {

	switch v := data.(type) {
	case wxbeacon2.WxEPData:
		fmt.Printf("DataType:EP ")
		fmt.Printf("Temp:%g ", v.Temp)
		fmt.Printf("Humid:%g ", v.Humid)
		fmt.Printf("AmbientLight:%d ", v.AmbientLight)
		fmt.Printf("UV Index:%g ", v.UVIndex)
		fmt.Printf("Pressure:%g ", v.Pressure)
		fmt.Printf("SoundNoise:%g ", v.SoundNoise)
		fmt.Printf("DisconfortIndex:%g ", v.DisconfortIndex)
		fmt.Printf("HeatStroke:%g ", v.HeatStroke)
		fmt.Printf("Battery:%gV\n", v.VBattery)
	case wxbeacon2.WxIMData:
		fmt.Printf("DataType:IM ")
		fmt.Printf("Temp:%g ", v.Temp)
		fmt.Printf("Humid:%g ", v.Humid)
		fmt.Printf("AmbientLight:%d ", v.AmbientLight)
		fmt.Printf("UV Index:%g ", v.UVIndex)
		fmt.Printf("Pressure:%g ", v.Pressure)
		fmt.Printf("SoundNoise:%g ", v.SoundNoise)
		fmt.Printf("Acceleration X:%d ", v.AccelerationX)
		fmt.Printf("Acceleration Y:%d ", v.AccelerationY)
		fmt.Printf("Acceleration Z:%d ", v.AccelerationZ)
		fmt.Printf("Battery:%gV\n", v.VBattery)
	}
}
