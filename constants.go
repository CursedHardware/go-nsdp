package nsdp

import (
	_ "embed"
	"encoding/hex"
)

type Version int

const (
	Version1 Version = 1
	Version2 Version = 2
)

type Command uint8

const (
	CommandReadRequest   Command = 0x01
	CommandReadResponse  Command = 0x02
	CommandWriteRequest  Command = 0x03
	CommandWriteResponse Command = 0x04
)

func (c Command) String() string {
	switch c {
	case CommandReadRequest:
		return "ReadRequest"
	case CommandReadResponse:
		return "ReadResponse"
	case CommandWriteRequest:
		return "WriteRequest"
	case CommandWriteResponse:
		return "WriteResponse"
	}
	return hex.EncodeToString([]byte{byte(c)})
}
