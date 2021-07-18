package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
)

const (
	MapObjectTypeEntity = iota
	MapObjectTypeBlock
)

// MapTrackedObject is an object on a map that is 'tracked' by the client, such as an entity or a block. This
// object may move, which is handled client-side.
type MapTrackedObject struct {
	// Type is the type of the tracked object. It is either MapObjectTypeEntity or MapObjectTypeBlock.
	Type int32
	// EntityUniqueID is the unique ID of the entity, if the tracked object was an entity. It needs not to be
	// filled out if Type is not MapObjectTypeEntity.
	EntityUniqueID int64
	// BlockPosition is the position of the block, if the tracked object was a block. It needs not to be
	// filled out if Type is not MapObjectTypeBlock.
	BlockPosition BlockPos
}

// MapDecoration is a fixed decoration on a map: Its position or other properties do not change automatically
// client-side.
type MapDecoration struct {
	// Type is the type of the map decoration. The type specifies the shape (and sometimes the colour) that
	// the map decoration gets.
	Type byte
	// Rotation is the rotation of the map decoration. It is byte due to the 16 fixed directions that the
	// map decoration may face.
	Rotation byte
	// X is the offset on the X axis in pixels of the decoration.
	X byte
	// Y is the offset on the Y axis in pixels of the decoration.
	Y byte
	// Label is the name of the map decoration. This name may be of any value.
	Label string
	// Colour is the colour of the map decoration. Some map decoration types have a specific colour set
	// automatically, whereas others may be changed.
	Colour color.RGBA
}

// MapTrackedObj reads a MapTrackedObject from buf into x.
func MapTrackedObj(buf *bytes.Buffer, x *MapTrackedObject) error {
	if err := binary.Read(buf, binary.LittleEndian, &x.Type); err != nil {
		return wrap(err)
	}
	switch x.Type {
	case MapObjectTypeEntity:
		return wrap(Varint64(buf, &x.EntityUniqueID))
	case MapObjectTypeBlock:
		return wrap(UBlockPosition(buf, &x.BlockPosition))
	default:
		return fmt.Errorf("unknown map tracked object type %v", x.Type)
	}
}

// WriteMapTrackedObj writes a MapTrackedObject xx to buf.
func WriteMapTrackedObj(buf *bytes.Buffer, x MapTrackedObject) error {
	if err := binary.Write(buf, binary.LittleEndian, x.Type); err != nil {
		return wrap(err)
	}
	switch x.Type {
	case MapObjectTypeEntity:
		return wrap(WriteVarint64(buf, x.EntityUniqueID))
	case MapObjectTypeBlock:
		return wrap(WriteUBlockPosition(buf, x.BlockPosition))
	default:
		panic(fmt.Sprintf("invalid map tracked object type %v", x.Type))
	}
}

// MapDeco reads a MapDecoration from buf into x.
func MapDeco(buf *bytes.Buffer, x *MapDecoration) error {
	return chainErr(
		binary.Read(buf, binary.LittleEndian, &x.Type),
		binary.Read(buf, binary.LittleEndian, &x.Rotation),
		binary.Read(buf, binary.LittleEndian, &x.X),
		binary.Read(buf, binary.LittleEndian, &x.Y),
		String(buf, &x.Label),
		VarRGBA(buf, &x.Colour),
	)
}

// WriteMapDeco writes a MapDecoration x to buf.
func WriteMapDeco(buf *bytes.Buffer, x MapDecoration) error {
	return chainErr(
		binary.Write(buf, binary.LittleEndian, x.Type),
		binary.Write(buf, binary.LittleEndian, x.Rotation),
		binary.Write(buf, binary.LittleEndian, x.X),
		binary.Write(buf, binary.LittleEndian, x.Y),
		WriteString(buf, x.Label),
		WriteVarRGBA(buf, x.Colour),
	)
}

// VarRGBA reads an RGBA value from buf into x packed into a varuint32.
func VarRGBA(buf *bytes.Buffer, x *color.RGBA) error {
	var v uint32
	err := wrap(Varuint32(buf, &v))
	*x = color.RGBA{
		R: byte(v),
		G: byte(v >> 8),
		B: byte(v >> 16),
		A: byte(v >> 24),
	}
	return err
}

// WriteVarRGBA writes an RGBA value to buf by packing it into a varuint32.
func WriteVarRGBA(buf *bytes.Buffer, x color.RGBA) error {
	return wrap(WriteVaruint32(buf, uint32(x.R)|uint32(x.G)<<8|uint32(x.B)<<16|uint32(x.A)<<24))
}
