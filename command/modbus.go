package command

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
	payload []byte
}

const modbusComAddr byte = 0xf7

type ModbusCommandType byte

const (
	ModbusCommandTypeRead       ModbusCommandType = 0x03
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

	return &ModbusCommand{payload: p}
}

func modbusChecksum(b []byte) uint16 {
	crc := uint16(0xffff)
	for _, v := range b {
		crc = (crc >> 8) ^ modbusCrcTable[(crc^uint16(v))&0xff]
	}
	return crc
}

func (cmd ModbusCommand) String() string { return string(cmd.payload) }

func (cmd ModbusCommand) validateResponse(p []byte) ([]byte, error) {
	return p[5 : len(p)-2], nil
}
