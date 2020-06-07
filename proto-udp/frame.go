package proto_udp

import (
	"encoding/binary"
	"fmt"
)

type Frame struct {
	ConnectionId uint32
	Opcode       uint8
	MessageId    uint8
	BlockId      uint16
	RawMessage   []byte
}

func frameDecode(byts []byte) (*Frame, error) {
	if len(byts) < 8 {
		return nil, fmt.Errorf("invalid length")
	}

	f := &Frame{}

	f.ConnectionId = binary.LittleEndian.Uint32(byts[0:4])
	f.Opcode = byts[4]
	f.MessageId = byts[5]
	f.BlockId = binary.LittleEndian.Uint16(byts[6:8])
	f.RawMessage = byts[8:]

	return f, nil
}

func frameEncode(f *Frame) ([]byte, error) {
	byts := make([]byte, 8+len(f.RawMessage))

	binary.LittleEndian.PutUint32(byts[:4], f.ConnectionId)
	byts[4] = f.Opcode
	byts[5] = f.MessageId
	binary.LittleEndian.PutUint16(byts[6:8], f.BlockId)
	copy(byts[8:], f.RawMessage)

	return byts, nil
}
