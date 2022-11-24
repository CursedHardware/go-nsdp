package report

import "net"

type ScannedReport struct {
	DeviceModel    string           `nsdp-scan:"0001"`
	DeviceName     string           `nsdp-scan:"0003"`
	MAC            net.HardwareAddr `nsdp-scan:"0004"`
	IP             net.IP           `nsdp-scan:"0006"`
	Mask           net.IP           `nsdp-scan:"0007"`
	Gateway        net.IP           `nsdp-scan:"0008"`
	UseDHCP        bool             `nsdp-scan:"000b"`
	ActiveFirmware uint8            `nsdp-scan:"000c"`
	Firmware1      string           `nsdp-scan:"000d"`
	Firmware2      string           `nsdp-scan:"000e"`
	Ports          uint8            `nsdp-scan:"6000"`
	SerialNumber   string           `nsdp-scan:"7800"`
}

func (r ScannedReport) IPNet() *net.IPNet {
	return &net.IPNet{IP: r.IP, Mask: net.IPMask(r.Mask)}
}
