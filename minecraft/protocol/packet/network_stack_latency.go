package packet

import (
	"bytes"
	"encoding/binary"
)

// NetworkStackLatency is sent by the server (and the client, on development builds) to measure the latency
// over the entire Minecraft stack, rather than the RakNet latency. It has other usages too, such as the
// ability to be used as some kind of acknowledgement packet, to know when the client has received a certain
// other packet.
type NetworkStackLatency struct {
	// Timestamp is the timestamp of the network stack latency packet. The client will, if NeedsResponse is
	// set to true, send a NetworkStackLatency packet with this same timestamp packet in response.
	Timestamp int64
	// NeedsResponse specifies if the sending side of this packet wants a response to the packet, meaning that
	// the other side should send a NetworkStackLatency packet back.
	NeedsResponse bool
}

// ID ...
func (*NetworkStackLatency) ID() uint32 {
	return IDNetworkStackLatency
}

// Marshal ...
func (pk *NetworkStackLatency) Marshal(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.LittleEndian, pk.Timestamp)
	_ = binary.Write(buf, binary.LittleEndian, pk.NeedsResponse)
}

// Unmarshal ...
func (pk *NetworkStackLatency) Unmarshal(buf *bytes.Buffer) error {
	return chainErr(
		binary.Read(buf, binary.LittleEndian, &pk.Timestamp),
		binary.Read(buf, binary.LittleEndian, &pk.NeedsResponse),
	)
}
