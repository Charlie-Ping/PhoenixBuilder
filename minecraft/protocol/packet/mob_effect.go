package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	MobEffectAdd = iota + 1
	MobEffectModify
	MobEffectRemove
)

const (
	EffectSpeed = iota + 1
	EffectSlowness
	EffectHaste
	EffectMiningFatigue
	EffectStrength
	EffectInstantHealth
	EffectInstantDamage
	EffectJumpBoost
	EffectNausea
	EffectRegeneration
	EffectResistance
	EffectFireResistance
	EffectWaterBreathing
	EffectInvisibility
	EffectBlindness
	EffectNightVision
	EffectHunger
	EffectWeakness
	EffectPoison
	EffectWither
	EffectHealthBoost
	EffectAbsorption
	EffectSaturation
	EffectLevitation
	EffectFatalPoison
	EffectConduitPower
)

// MobEffect is sent by the server to apply an effect to the player, for example an effect like poison. It may
// also be used to modify existing effects, or removing them completely.
type MobEffect struct {
	// EntityRuntimeID is the runtime ID of the entity. The runtime ID is unique for each world session, and
	// entities are generally identified in packets using this runtime ID.
	EntityRuntimeID uint64
	// Operation is the operation of the packet. It is either MobEffectAdd, MobEffectModify or MobEffectRemove
	// and specifies the result of the packet client-side.
	Operation byte
	// EffectType is the ID of the effect to be added, removed or modified. It is one of the constants that
	// may be found above.
	EffectType int32
	// Amplifier is the amplifier of the effect. Take note that the amplifier is not the same as the effect's
	// level. The level is usually one higher than the amplifier, and the amplifier can actually be negative
	// to reverse the behaviour effect.
	Amplifier int32
	// Particles specifies if viewers of the entity that gets the effect shows particles around it. If set to
	// false, no particles are emitted around the entity.
	Particles bool
	// Duration is the duration of the effect in seconds. After the duration has elapsed, the effect will be
	// removed automatically client-side.
	Duration int32
}

// ID ...
func (*MobEffect) ID() uint32 {
	return IDMobEffect
}

// Marshal ...
func (pk *MobEffect) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteVaruint64(buf, pk.EntityRuntimeID)
	_ = binary.Write(buf, binary.LittleEndian, pk.Operation)
	_ = protocol.WriteVarint32(buf, pk.EffectType)
	_ = protocol.WriteVarint32(buf, pk.Amplifier)
	_ = binary.Write(buf, binary.LittleEndian, pk.Particles)
	_ = protocol.WriteVarint32(buf, pk.Duration)
}

// Unmarshal ...
func (pk *MobEffect) Unmarshal(buf *bytes.Buffer) error {
	return chainErr(
		protocol.Varuint64(buf, &pk.EntityRuntimeID),
		binary.Read(buf, binary.LittleEndian, &pk.Operation),
		protocol.Varint32(buf, &pk.EffectType),
		protocol.Varint32(buf, &pk.Amplifier),
		binary.Read(buf, binary.LittleEndian, &pk.Particles),
		protocol.Varint32(buf, &pk.Duration),
	)
}
