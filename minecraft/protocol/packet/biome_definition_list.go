package packet

import (
	"bytes"
	"math"
)

// BiomeDefinitionList is sent by the server to let the client know all biomes that are available and
// implemented on the server side. It is much like the AvailableActorIdentifiers packet, but instead
// functions for biomes.
type BiomeDefinitionList struct {
	// SerialisedBiomeDefinitions is a network NBT serialised compound of all definitions of biomes that are
	// available on the server.
	SerialisedBiomeDefinitions []byte
}

// ID ...
func (*BiomeDefinitionList) ID() uint32 {
	return IDBiomeDefinitionList
}

// Marshal ...
func (pk *BiomeDefinitionList) Marshal(buf *bytes.Buffer) {
	_, _ = buf.Write(pk.SerialisedBiomeDefinitions)
}

// Unmarshal ...
func (pk *BiomeDefinitionList) Unmarshal(buf *bytes.Buffer) error {
	pk.SerialisedBiomeDefinitions = buf.Next(math.MaxInt32)
	return nil
}
