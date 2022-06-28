package report

import "net"

type ScannedReport struct {
	DeviceModel    string           `nsdp:"0001"`
	DeviceName     string           `nsdp:"0003"`
	MAC            net.HardwareAddr `nsdp:"0004"`
	IP             net.IP           `nsdp:"0006"`
	Mask           net.IP           `nsdp:"0007"`
	Gateway        net.IP           `nsdp:"0008"`
	UseDHCP        bool             `nsdp:"000b"`
	ActiveFirmware uint8            `nsdp:"000c"`
	Firmware1      string           `nsdp:"000d"`
	Firmware2      string           `nsdp:"000e"`
	Ports          uint8            `nsdp:"6000"`
	SerialNumber   string           `nsdp:"7800"`
}

func (r ScannedReport) IPNet() *net.IPNet {
	return &net.IPNet{IP: r.IP, Mask: net.IPMask(r.Mask)}
}
