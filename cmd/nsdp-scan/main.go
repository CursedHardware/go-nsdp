package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/CursedHardware/go-nsdp"
	. "github.com/CursedHardware/go-nsdp/report"
)

func main() {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		return
	}
	client, err := nsdp.NewClient(iface.HardwareAddr, net.IP{192, 168, 1, 145}, nsdp.Version2)
	if err != nil {
		return
	}
	onScanning(client)
}

func onScanning(client *nsdp.Client) {
	writer := csv.NewWriter(os.Stdout)
	_ = writer.Write([]string{
		"Model",
		"Name",
		"IP",
		"Active Firmware",
		"Firmware (slot 1)",
		"Firmware (slot 2)",
		"Ports",
		"Serial Number",
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	reports := make(chan *ScannedReport)
	_ = Scan(ctx, client, reports)
	for report := range reports {
		_ = writer.Write([]string{
			report.DeviceModel,
			report.DeviceName,
			report.IPNet().String(),
			fmt.Sprint(report.ActiveFirmware),
			report.Firmware1,
			report.Firmware2,
			fmt.Sprint(report.Ports),
			report.SerialNumber,
		})
		writer.Flush()
	}
}
