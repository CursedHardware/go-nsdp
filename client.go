package nsdp

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net"
	"time"
)

type Callback func(*Message)

type Client struct {
	conn        net.PacketConn
	destination net.Addr
	managerId   net.HardwareAddr
	rand        *rand.Rand
	callbacks   map[uint32]Callback
	scanning    map[uint32]Callback
}

func NewClient(managerId net.HardwareAddr, localIP net.IP, version Version) (client *Client, err error) {
	source := &net.UDPAddr{IP: localIP}
	destination := &net.UDPAddr{IP: net.IPv4bcast}
	switch version {
	case Version1:
		source.Port, destination.Port = 63323, 63324
	case Version2:
		source.Port, destination.Port = 63321, 63322
	}
	client = &Client{
		managerId:   managerId,
		destination: destination,
		rand:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	client.conn, err = net.ListenUDP("udp4", source)
	client.callbacks = make(map[uint32]Callback)
	client.scanning = make(map[uint32]Callback)
	go client.watch()
	return
}

func (c *Client) watch() {
	packet := make([]byte, 0x400)
	var message *Message
	for {
		if n, _, err := c.conn.ReadFrom(packet); err != nil {
			continue
		} else {
			message = new(Message)
			_, _ = message.ReadFrom(bytes.NewReader(packet[:n]))
		}
		if callback, ok := c.callbacks[message.id()]; ok {
			go callback(message)
		} else if len(c.scanning) > 0 {
			for _, event := range c.scanning {
				go event(message)
			}
		}
	}
}

// Scan Design using method:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	messages := make(chan *nsdp.Message)
//	go client.Scan(ctx, nsdp.ScanTags(), messages)
//	for message := range messages {
//	    // your business
//	}
func (c *Client) Scan(context context.Context, tags Tags, onCallback Callback) (err error) {
	id := c.rand.Uint32()
	c.scanning[id] = onCallback
	defer delete(c.scanning, id)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-context.Done():
			return
		case <-ticker.C:
			request := &Message{
				Header: c.Header(CommandReadRequest, nil),
				Tags:   tags,
			}
			if err = c.sendRequest(request); err != nil {
				err = errors.New("nsdp: send request failed")
				return
			}
		}
	}
}

func (c *Client) Request(context context.Context, request *Message) (response *Message, err error) {
	returns := make(chan *Message, 1)
	c.callbacks[request.id()] = func(message *Message) { returns <- message }
	defer delete(c.callbacks, request.id())
	if err = c.sendRequest(request); err != nil {
		err = errors.New("nsdp: send request failed")
		return
	}
	select {
	case <-context.Done():
		err = errors.New("nsdp: failed to wait response")
	case response = <-returns:
	}
	return
}

func (c *Client) sendRequest(request *Message) (err error) {
	var buf bytes.Buffer
	_, err = request.WriteTo(&buf)
	if err == nil {
		_, err = c.conn.WriteTo(buf.Bytes(), c.destination)
	}
	return
}

func (c *Client) Header(command Command, agentID net.HardwareAddr) Header {
	header := Header{
		Version:   1,
		Command:   command,
		Signature: [4]byte{'N', 'S', 'D', 'P'},
	}
	_, _ = c.rand.Read(header.Sequence[:])
	copy(header.ManagerID[:], c.managerId)
	copy(header.AgentID[:], agentID)
	return header
}

func (c *Client) Close() error {
	return c.conn.Close()
}

type DeviceClient struct {
	*Client
	AgentID  net.HardwareAddr
	Password []byte
}

func (c *DeviceClient) Read(context context.Context, tags Tags) (*Message, error) {
	request := &Message{
		Header: c.Client.Header(CommandReadRequest, c.AgentID),
		Tags:   tags,
	}
	return c.Client.Request(context, request)
}

func (c *DeviceClient) Write(context context.Context, tags Tags) (*Message, error) {
	request := &Message{
		Header: c.Client.Header(CommandWriteRequest, c.AgentID),
		Tags:   tags,
	}
	return c.Client.Request(context, request)
}

func (c *DeviceClient) Set(context context.Context, tags Tags) (response *Message, err error) {
	var salt *Message
	if salt, err = c.Read(context, Tags{0x0017: nil}); err != nil {
		return
	}
	tags[0x001a] = AuthV2Password(c.Password, salt.AgentID[:], salt.Tags[0x0017])
	return c.Write(context, tags)
}

// AuthV2Password
// based https://github.com/yaamai/go-nsdp/blob/d54b436f/nsdp/auth_v2.go modify
func AuthV2Password(passphrase, mac, salt []byte) []byte {
	var key [20]byte
	copy(key[:], passphrase)
	return []byte{
		salt[3] ^ salt[2] ^ mac[1] ^ mac[5] ^ key[0] ^ key[1] ^ key[2],
		salt[3] ^ salt[1] ^ mac[4] ^ mac[0] ^ key[3] ^ key[4] ^ key[5],
		salt[0] ^ salt[2] ^ mac[3] ^ mac[2] ^ key[6] ^ key[7] ^ key[8],
		salt[0] ^ salt[1] ^ mac[4] ^ mac[5] ^ key[9] ^ key[10] ^ key[11],
		salt[3] ^ salt[2] ^ mac[1] ^ mac[5] ^ key[12] ^ key[13] ^ key[14],
		salt[3] ^ salt[1] ^ mac[4] ^ mac[0] ^ key[15] ^ key[16] ^ key[17],
		salt[0] ^ salt[2] ^ mac[3] ^ mac[2] ^ key[18] ^ key[19] ^ key[0],
		salt[0] ^ salt[1] ^ mac[4] ^ mac[5] ^ key[1] ^ key[3] ^ key[5],
	}
}
