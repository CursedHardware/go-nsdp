package nsdp

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"io"
	"sort"
)

type Header struct {
	Version    byte
	Command    Command
	Status     uint16
	FailureTLV [4]byte
	ManagerID  [6]byte
	AgentID    [6]byte
	Sequence   [4]byte
	Signature  [4]byte
}

func (m *Header) id() uint32 {
	return binary.BigEndian.Uint32(m.Sequence[:])
}

type Tag uint16

func (t Tag) String() string {
	var tag [2]byte
	binary.BigEndian.PutUint16(tag[:], uint16(t))
	return hex.EncodeToString(tag[:])
}

type Tags map[Tag][]byte

func (t Tags) WriteTo(w io.Writer) (n int64, err error) {
	keys := make([]int, 0)
	for tag := range t {
		keys = append(keys, int(tag))
	}
	sort.Ints(keys)
	var buf bytes.Buffer
	for _, key := range keys {
		block := t[Tag(key)]
		_ = binary.Write(&buf, binary.BigEndian, uint16(key))
		_ = binary.Write(&buf, binary.BigEndian, uint16(len(block)))
		buf.Write(block)
	}
	return buf.WriteTo(w)
}

func (t *Tags) ReadFrom(r io.Reader) (n int64, err error) {
	*t = make(map[Tag][]byte)
	var tag Tag
	var length uint16
	for err != io.EOF {
		err = binary.Read(r, binary.BigEndian, &tag)
		if err == nil {
			err = binary.Read(r, binary.BigEndian, &length)
		}
		if err == nil {
			block := make([]byte, length)
			_, err = r.Read(block)
			(*t)[tag] = block
		}
	}
	return
}

type Message struct {
	Header
	Tags Tags
}

func (m *Message) WriteTo(w io.Writer) (n int64, err error) {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.BigEndian, &m.Header)
	m.Tags[0x0000] = nil // start marker
	m.Tags[0xffff] = nil // end marker
	_, _ = m.Tags.WriteTo(&buf)
	delete(m.Tags, 0x0000)
	delete(m.Tags, 0xffff)
	return buf.WriteTo(w)
}

func (m *Message) ReadFrom(r io.Reader) (n int64, err error) {
	err = binary.Read(r, binary.BigEndian, &m.Header)
	if err == nil {
		_, err = m.Tags.ReadFrom(r)
	}
	delete(m.Tags, 0x0000) // start marker
	delete(m.Tags, 0xffff) // end marker
	return
}
