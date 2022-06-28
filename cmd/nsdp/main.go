package main

import (
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/CursedHardware/go-nsdp"
	. "github.com/CursedHardware/go-nsdp/report"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	iface, err := net.InterfaceByName("en0")
	if err != nil {
		return
	}
	client, err := nsdp.NewClient(iface.HardwareAddr, net.IP{192, 168, 1, 127}, nsdp.Version2)
	if err != nil {
		return
	}
	onPasswordSalt(&nsdp.DeviceClient{
		Client: client,
	})
	//onScanning(client)
}

func onPasswordSalt(client *nsdp.DeviceClient) {
	message, err := client.Read(context.Background(), nsdp.Tags{0x0017: nil})
	if err != nil {
		return
	}
	log.Println(
		hex.EncodeToString(message.AgentID[:]),
		hex.EncodeToString(message.Tags[0x0017]),
	)
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	founded := make(map[string]bool)
	messages := make(chan *nsdp.Message)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	go client.Scan(ctx, nsdp.ScanTags(), messages)
	for message := range messages {
		select {
		case <-signals:
			cancel()
			return
		default:
		}
		report := new(ScannedReport)
		_ = UnmarshalReport(message, report)
		if _, ok := founded[report.SerialNumber]; ok {
			continue
		}
		founded[report.SerialNumber] = true
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
