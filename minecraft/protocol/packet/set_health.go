package packet

import (
	"bytes"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// SetHealth is sent by the server. It sets the health of the player it is sent to. The SetHealth packet
// should no longer be used. Instead, the health attribute should be used so that the health and maximum
// health may be changed directly.
type SetHealth struct {
	// Health is the new health of the player.
	Health int32
}

// ID ...
func (*SetHealth) ID() uint32 {
	return IDSetHealth
}

// Marshal ...
func (pk *SetHealth) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteVarint32(buf, pk.Health)
}

// Unmarshal ...
func (pk *SetHealth) Unmarshal(buf *bytes.Buffer) error {
	return protocol.Varint32(buf, &pk.Health)
}
