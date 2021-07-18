package packet

import (
	"bytes"
	"encoding/binary"
)

// ClientCacheStatus is sent by the client to the server at the start of the game. It is sent to let the
// server know if it supports the client-side blob cache. Clients such as Nintendo Switch do not support the
// cache, and attempting to use it anyway will fail.
type ClientCacheStatus struct {
	// Enabled specifies if the blob cache is enabled. If false, the server should not attempt to use the
	// blob cache. If true, it may do so, but it may also choose not to use it.
	Enabled bool
}

// ID ...
func (pk *ClientCacheStatus) ID() uint32 {
	return IDClientCacheStatus
}

// Marshal ...
func (pk *ClientCacheStatus) Marshal(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.LittleEndian, pk.Enabled)
}

// Unmarshal ...
func (pk *ClientCacheStatus) Unmarshal(buf *bytes.Buffer) error {
	return binary.Read(buf, binary.LittleEndian, &pk.Enabled)
}
