package nsdp

import "errors"

var (
	ErrSendRequestFailed    = errors.New("nsdp-scan: send request failed")
	ErrFailedToWaitResponse = errors.New("nsdp-scan: failed to wait response")
)
