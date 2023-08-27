package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/walkure/go-wxbeacon2"
)

func main() {
	wxbeacon2.SetLogger(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	dev := wxbeacon2.NewDevice("ZZ:ZZ:ZZ:ZZ:ZZ:ZZ", printData)
	err := dev.WaitForReceiveData()
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to open device:%v", err))
		return
	}

	select {}
}

func printData(data interface{}) {
	v, ok := data.(slog.LogValuer)
	if ok {
		slog.Info("packet received", "data", v)
	} else {
		slog.Error(fmt.Sprintf("Unknown data type(%T)", data))
	}
}
