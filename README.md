# WxBeacon2(2JCIE-BL01) Go Receiver

Go Receiver for [WxBeacon2](https://weathernews.jp/smart/wxbeacon2/)([2JCIE-BL](https://www.omron.co.jp/ecb/product-detail?partNumber=2JCIE-BL)01).

# dependencies

- [gatt](https://github.com/walkure/gatt)
  - modified from [bettercap/gatt](https://github.com/bettercap/gatt) which is most active repository forked from [paypal/gatt](https://github.com/paypal/gatt) in May 2021.


# references

 - 2JCIE-BL01 Communication Interface Manual 
   - [Japanese PDF](https://omronfs.omron.com/ja_JP/ecb/products/pdf/CDSC-015.pdf)
   - [English HTML](https://omronmicrodevices.github.io/products/2jcie-bl01/communication_if_manual.html)

# sample code

see `cmd/main.go`

# notice
 - Currently this code supports `Limited/General Broadcaster` mode(see reference).
 - You should run this code as `root` user(or CAP_NET_ADMIN capability).
