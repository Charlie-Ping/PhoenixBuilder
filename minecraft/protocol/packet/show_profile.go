package packet

import (
	"bytes"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// ShowProfile is sent by the server to show the XBOX Live profile of one player to another.
type ShowProfile struct {
	// XUID is the XBOX Live User ID of the player whose profile should be shown to the player. If it is not
	// a valid XUID, the client ignores the packet.
	XUID string
}

// ID ...
func (*ShowProfile) ID() uint32 {
	return IDShowProfile
}

// Marshal ...
func (pk *ShowProfile) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteString(buf, pk.XUID)
}

// Unmarshal ...
func (pk *ShowProfile) Unmarshal(buf *bytes.Buffer) error {
	return protocol.String(buf, &pk.XUID)
}
