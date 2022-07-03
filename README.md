# WxBeacon2(2JCIE-BL01) Go Receiver

Go Receiver for [WxBeacon2](https://weathernews.jp/smart/wxbeacon2/)([2JCIE-BL](https://www.omron.co.jp/ecb/product-detail?partNumber=2JCIE-BL)01).

# dependencies

- [bettercap/gatt](https://github.com/bettercap/gatt)
  - most active repository forked from [paypal/gatt](https://github.com/paypal/gatt) in May 2021.


# references

 - 2JCIE-BL01 Communication Interface Manual 
   - [Japanese PDF](https://omronfs.omron.com/ja_JP/ecb/products/pdf/CDSC-015.pdf)
   - [English HTML](https://omronmicrodevices.github.io/products/2jcie-bl01/communication_if_manual.html)

# sample code
```go
func main() {
	dev := wxbeacon2.NewDevice("ZZ:ZZ:ZZ:ZZ:ZZ:ZZ", process)
	err := dev.WaitForReceiveData()
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	select {}
}
func process(data interface{}){
	switch v := data.(type) {
	case wxbeacon2.WxIMData:
        // process IM Mode Data
	case wxbeacon2.WxEPData:
        // process EP Mode Data
}
```

# notice
 - Currently this code supports `Limited/General Broadcaster` mode(see reference).
 - You should run this code as `root` user(or CAP_NET_ADMIN capability).
