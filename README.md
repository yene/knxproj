# knxproj - KNX project parser

I ported an existing project parser from Java to Go, which resulted in a bit of an ugly code.

Export used devices as JSON.
`go run cmd/main.go Project.knxproj`

## TODO
- [ ] Test with two level group addresses.
- [ ] Test against all knxproj files.
- [ ] Sort addresses into a nice tree.
- [ ] When converting address to string, respect the configured format.
- [ ] Handle IsActive, somehow find out if device needs a full reconfiguration.
- [ ] Improve json structs http://attilaolah.eu/2014/09/10/json-and-struct-composition-in-go/
- [ ] Convert bit size without DPT, to a good dpt. 1bit = dpt 1. 2bit = dpt2, 4bit = 3.007
- [ ] Study this project https://github.com/owagner/knx2mqtt

## Resources
* [tuxedo0801/ETS4Reader](https://github.com/tuxedo0801/ETS4Reader)
* [xmljson2struct](https://github.com/wicast/xj2s)
* [migrating OOP to Go](https://github.com/luciotato/golang-notes/blob/master/OOP.md)
* [java implementation](https://github.com/IOT-DSA/dslink-java-knx/tree/master/src/main/java/org/dsa/iot/knx)
* [Custom unmarshal](http://stackoverflow.com/a/30857066/279890)
