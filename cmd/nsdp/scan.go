package main

import (
	"context"
	"github.com/CursedHardware/go-nsdp"
	. "github.com/CursedHardware/go-nsdp/report"
)

func Scan(context context.Context, client *nsdp.Client, reports chan *ScannedReport) error {
	founded := make(map[string]bool)
	return client.Scan(context, nsdp.ScanTags(), func(message *nsdp.Message) {
		report := new(ScannedReport)
		_ = UnmarshalReport(message, report)
		if _, ok := founded[report.SerialNumber]; !ok {
			founded[report.SerialNumber] = true
			reports <- report
		}
	})
}
