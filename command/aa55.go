package command

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

const (
	aa55Header               = "AA55C07F"
	aa55ResponseLengthIndex  = 6
	aa55ResponseLengthOffset = 9
)

type AA55Command struct {
	payload      []byte
	responseType string
}

func NewAA55(payload, responseType string) (*AA55Command, error) {
	bytes, err := hex.DecodeString(aa55Header + payload)
	if err != nil {
		return nil, fmt.Errorf("error parsing payload: %s", err)
	}

	bytes = append(bytes, aa55Checksum(bytes)...)

	return &AA55Command{payload: bytes, responseType: responseType}, nil
}

func aa55Checksum(payload []byte) []byte {
	var v uint16
	for _, b := range payload {
		v += uint16(b)
	}

	c := make([]byte, 4)
	binary.BigEndian.PutUint16(c, v)
	return c
}

func (cmd AA55Command) String() string { return string(cmd.payload) }

func (cmd AA55Command) ValidateResponse(p []byte) ([]byte, error) {
	if len(p) < 8 {
		return nil, fmt.Errorf("response truncated")
	}

	expectedLen := int(p[aa55ResponseLengthIndex] + aa55ResponseLengthOffset)
	if len(p) != expectedLen {
		return nil, fmt.Errorf("unexpected response length %d (expected %d)", len(p), expectedLen)
	}

	responseType := hex.EncodeToString(p[4:6])
	if responseType != cmd.responseType {
		return nil, fmt.Errorf("unexpected response type `%s` (expected `%s`)", responseType, cmd.responseType)
	}

	var sum uint16
	for _, b := range p[:len(p)-2] {
		sum += uint16(b)
	}
	expSum := binary.BigEndian.Uint16(p[len(p)-2:])
	if sum != expSum {
		return nil, fmt.Errorf("invalid response checksum %d (expected %d)", sum, expSum)
	}

	// FIXME: use correct offsets
	return p[5 : len(p)-2], nil
}
