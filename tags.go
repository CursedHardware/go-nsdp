package nsdp

func ScanTags() Tags {
	return Tags{
		0x0001: nil, // device model
		0x0003: nil, // device name
		0x0004: nil, // mac address
		0x0005: nil, // location
		0x0006: nil, // ip address
		0x0007: nil, // ip mask
		0x0008: nil, // gateway
		0x000b: nil, // use dhcp
		0x000c: nil, // active firmware
		0x000d: nil, // firmware (slot 1)
		0x000e: nil, // firmware (slot 2)
		0x6000: nil, // number of ports
		0x7800: nil, // serial number
	}
}
