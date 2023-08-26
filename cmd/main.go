package main

import (
	"fmt"
	"github.com/walkure/go-wxbeacon2"
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
	v, ok := data.(fmt.Stringer)
	if ok {
		log.Printf(v.String())
	} else {
		log.Printf("Unknown data type(%T)", v)
	}
}
