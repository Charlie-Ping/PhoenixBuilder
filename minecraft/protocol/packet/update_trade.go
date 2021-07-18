package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// UpdateTrade is sent by the server to update the trades offered by a villager to a player. It is sent at the
// moment that a player interacts with a villager.
type UpdateTrade struct {
	// WindowID is the ID that identifies the trading window that the client currently has opened.
	WindowID byte
	// WindowType is an identifier specifying the type of the window opened. In vanilla, it appears this is
	// always filled out with 15.
	WindowType byte
	// UnknownInt32 ...
	UnknownInt32 int32
	// TradeTier is the tier of the villager that the player is trading with. The tier starts at 0 with a
	// first two offers being available, after which two additional offers are unlocked each time the tier
	// becomes one higher.
	TradeTier int32
	// VillagerUniqueID is the unique ID of the villager entity that the player is trading with. The
	// TradeTier sent above applies to this villager.
	VillagerUniqueID int64
	// EntityUniqueID is the unique ID of the entity (usually a player) for which the trades are updated. The
	// updated trades may apply only to this entity.
	EntityUniqueID int64
	// DisplayName is the name displayed at the top of the trading UI. It is usually used to represent the
	// profession of the villager in the UI.
	DisplayName string
	// NewTradeUI specifies if the villager should be using the new trade UI (The one added in 1.11.) rather
	// than the old one. This should usually be set to true.
	NewTradeUI bool
	// DemandBasedPrices specifies if the prices of the villager's offers are modified by an increase in
	// demand for the item. (A mechanic added in 1.11.) Buying more of the same item will increase the price
	// of that particular item.
	DemandBasedPrices bool
	// SerialisedOffers is a network NBT serialised compound of offers that the villager has.
	SerialisedOffers []byte
}

// ID ...
func (*UpdateTrade) ID() uint32 {
	return IDUpdateTrade
}

// Marshal ...
func (pk *UpdateTrade) Marshal(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.LittleEndian, pk.WindowID)
	_ = binary.Write(buf, binary.LittleEndian, pk.WindowType)
	_ = protocol.WriteVarint32(buf, pk.UnknownInt32)
	_ = protocol.WriteVarint32(buf, pk.TradeTier)
	_ = protocol.WriteVarint64(buf, pk.VillagerUniqueID)
	_ = protocol.WriteVarint64(buf, pk.EntityUniqueID)
	_ = protocol.WriteString(buf, pk.DisplayName)
	_ = binary.Write(buf, binary.LittleEndian, pk.NewTradeUI)
	_ = binary.Write(buf, binary.LittleEndian, pk.DemandBasedPrices)
	_, _ = buf.Write(pk.SerialisedOffers)
}

// Unmarshal ...
func (pk *UpdateTrade) Unmarshal(buf *bytes.Buffer) error {
	if err := chainErr(
		binary.Read(buf, binary.LittleEndian, &pk.WindowID),
		binary.Read(buf, binary.LittleEndian, &pk.WindowType),
		protocol.Varint32(buf, &pk.UnknownInt32),
		protocol.Varint32(buf, &pk.TradeTier),
		protocol.Varint64(buf, &pk.VillagerUniqueID),
		protocol.Varint64(buf, &pk.EntityUniqueID),
		protocol.String(buf, &pk.DisplayName),
		binary.Read(buf, binary.LittleEndian, &pk.NewTradeUI),
		binary.Read(buf, binary.LittleEndian, &pk.DemandBasedPrices),
	); err != nil {
		return err
	}
	pk.SerialisedOffers = make([]byte, buf.Len())
	_, err := buf.Read(pk.SerialisedOffers)
	return err
}
