package nsdp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthV2Password(t *testing.T) {
	expected := []byte{0xc4, 0xaf, 0x7c, 0x00, 0xa6, 0xc4, 0x1a, 0x7d}
	result := AuthV2Password(
		[]byte("password"),                         // password
		[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc}, // switch mac address
		[]byte{0x12, 0x34, 0x56, 0x78},             // salt
	)
	assert.Equal(t, expected, result)
}
