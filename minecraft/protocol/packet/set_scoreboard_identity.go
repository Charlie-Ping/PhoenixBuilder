package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	ScoreboardIdentityActionRegister = iota
	ScoreboardIdentityActionClear
)

// SetScoreboardIdentity is sent by the server to change the identity type of one of the entries on a
// scoreboard. This is used to change, for example, an entry pointing to a player, to a fake player when it
// leaves the server, and to change it back to a real player when it joins again.
// In non-vanilla situations, the packet is quite useless.
type SetScoreboardIdentity struct {
	// ActionType is the type of the action to execute. The action is either ScoreboardIdentityActionRegister
	// to associate an identity with the entry, or ScoreboardIdentityActionClear to remove associations with
	// an entity.
	ActionType byte
	// Entries is a list of all entries in the packet. Each of these entries points to one of the entries on
	// a scoreboard. Depending on ActionType, their identity will either be registered or cleared.
	Entries []protocol.ScoreboardIdentityEntry
}

// ID ...
func (*SetScoreboardIdentity) ID() uint32 {
	return IDSetScoreboardIdentity
}

// Marshal ...
func (pk *SetScoreboardIdentity) Marshal(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.LittleEndian, pk.ActionType)
	_ = protocol.WriteVaruint32(buf, uint32(len(pk.Entries)))
	for _, entry := range pk.Entries {
		_ = protocol.WriteVarint64(buf, entry.EntryID)
		if pk.ActionType == ScoreboardIdentityActionRegister {
			_ = protocol.WriteVarint64(buf, entry.EntityUniqueID)
		}
	}
}

// Unmarshal ...
func (pk *SetScoreboardIdentity) Unmarshal(buf *bytes.Buffer) error {
	var count uint32
	if err := chainErr(
		binary.Read(buf, binary.LittleEndian, &pk.ActionType),
		protocol.Varuint32(buf, &count),
	); err != nil {
		return err
	}
	pk.Entries = make([]protocol.ScoreboardIdentityEntry, count)
	for i := uint32(0); i < count; i++ {
		if err := protocol.Varint64(buf, &pk.Entries[i].EntryID); err != nil {
			return err
		}
		if pk.ActionType == ScoreboardIdentityActionRegister {
			if err := protocol.Varint64(buf, &pk.Entries[i].EntityUniqueID); err != nil {
				return err
			}
		}
	}
	return nil
}
