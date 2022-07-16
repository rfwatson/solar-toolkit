package command

//go:generate stringer -type=FailureCode

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type FailureCode byte

const (
	FailureCodeIllegalFunction FailureCode = iota + 1
	FailureCodeIllegalDataAddress
	FailureCodeIllegalDataValue
	FailureCodeSlaveDeviceFailure
	FailureCodeAcknowledge
	FailureCodeSlaveDeviceBusy
	FailureCodeNegativeAcknowledgement
	FailureCodeMemoryParityError
	FailureCodeGatewayPathUnavailable
	FailureCodeGatewayTargetDeviceFailedToRespond
)

var modbusCrcTable []uint16

func init() {
	for i := 0; i < 256; i++ {
		buffer := uint16(i << 1)
		var crc uint16
		for j := 8; j > 0; j-- {
			buffer >>= 1
			if (buffer^crc)&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
		modbusCrcTable = append(modbusCrcTable, crc)
	}
}

type ModbusCommand struct {
	payload     []byte
	commandType ModbusCommandType
	offset      uint16
	value       uint16
}

const modbusComAddr byte = 0xf7

type ModbusCommandType byte

const (
	ModbusCommandTypeRead ModbusCommandType = 0x03
	// TODO: implement write commands.
	ModbusCommandTypeWrite      ModbusCommandType = 0x06
	ModbusCommandTypeWriteMulti ModbusCommandType = 0x10
)

func NewModbus(commandType ModbusCommandType, offset uint16, value uint16) *ModbusCommand {
	var p []byte
	p = append(p, modbusComAddr)
	p = append(p, byte(commandType))
	p = append(p, byte((offset>>8)&0xff))
	p = append(p, byte(offset&0xff))
	p = append(p, byte((value>>8)&0xff))
	p = append(p, byte(value&0xff))

	sum := modbusChecksum(p)
	p = append(p, byte(sum&0xff))
	p = append(p, byte((sum>>8)&0xff))

	return &ModbusCommand{
		payload:     p,
		commandType: commandType,
		offset:      offset,
		value:       value,
	}
}

func modbusChecksum(b []byte) uint16 {
	crc := uint16(0xffff)
	for _, v := range b {
		crc = (crc >> 8) ^ modbusCrcTable[(crc^uint16(v))&0xff]
	}
	return crc
}

func (cmd ModbusCommand) String() string { return string(cmd.payload) }

// ValidateResponse validates the entire response and if valid returns the
// response body.
func (cmd ModbusCommand) ValidateResponse(p []byte) ([]byte, error) {
	if len(p) < 4 {
		return nil, errors.New("invalid response: response too short")
	}

	var expectedLen int
	cmdType := ModbusCommandType(p[3])

	switch cmdType {
	case ModbusCommandTypeRead:
		if uint16(p[4]) != cmd.value*2 {
			return nil, fmt.Errorf("short response: expected %d, got %d", p[4], cmd.value*2)
		}
		expectedLen = int(p[4]) + 7
		if len(p) < expectedLen {
			return nil, fmt.Errorf("invalid read length: expected %d, got %d", expectedLen, len(p))
		}
	case ModbusCommandTypeWrite, ModbusCommandTypeWriteMulti:
		panic("unimplemented")
	default:
		expectedLen = len(p)
	}

	offset := expectedLen - 2
	wantSum := modbusChecksum(p[2:offset])
	gotSum := binary.LittleEndian.Uint16(p[offset:])
	if wantSum != gotSum {
		return nil, fmt.Errorf("invalid CRC-16: want `%X`, got `%X`", wantSum, gotSum)
	}

	if p[3] != byte(cmd.commandType) {
		failureCode := FailureCode(p[4])
		return nil, fmt.Errorf("command failed with code: %d, error: %s", failureCode, failureCode.String())
	}

	return p[5:offset], nil
}
