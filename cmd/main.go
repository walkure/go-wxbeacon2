package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"

	"github.com/walkure/gatt"
	"github.com/walkure/go-wxbeacon2"
)

func main() {
	// Passive scanning
	d, err := gatt.NewDevice(gatt.LnxSetScanMode(false))

	if err != nil {
		panic(err)
	}

	var mu sync.Mutex
	seqno := uint8(0)

	d.Handle(gatt.PeripheralDiscovered(
		wxbeacon2.HandleWxBeacon2("", // "aa:zz:pp:ff:dd:cc"
			func(d wxbeacon2.WxData) {
				// GATT lib. calls this callback function with a new goroutine.
				mu.Lock()
				defer mu.Unlock()

				if seqno != d.WxDataSequence() {
					v, ok := d.(slog.LogValuer)
					if ok {
						slog.Info("packet received", "data", v)
					} else {
						slog.Error(fmt.Sprintf("Unknown data type(%T)", d))
					}

					seqno = d.WxDataSequence()
				}

			}, nil)))
	d.Init(onStateChanged)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	fmt.Println("interrupted. Bye~")

	d.StopScanning()
	fmt.Printf("scan stopped: %+v\n", d.Stop())
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
