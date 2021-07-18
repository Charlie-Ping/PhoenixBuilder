package packet

import (
	"bytes"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SetTime is sent by the server to update the current time client-side. The client actually advances time
// client-side by itself, so this packet does not need to be sent each tick. It is merely a means of
// synchronising time between server and client.
type SetTime struct {
	// Time is the current time. The time is not limited to 24000 (time of day), but continues progressing
	// after that.
	Time int32
}

// ID ...
func (*SetTime) ID() uint32 {
	return IDSetTime
}

// Marshal ...
func (pk *SetTime) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteVarint32(buf, pk.Time)
}

// Unmarshal ...
func (pk *SetTime) Unmarshal(buf *bytes.Buffer) error {
	return protocol.Varint32(buf, &pk.Time)
}
