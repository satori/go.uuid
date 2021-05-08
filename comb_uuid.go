package uuid

import (
	"bytes"
	"encoding/binary"
	"time"
)

type CombUUID UUID

// Equal returns true if this uuid and other uuid equals, otherwise false
func (u CombUUID) Equal(other CombUUID) bool {
	return bytes.Equal(u[:], other[:])
}

// Bytes returns bytes slice representation of UUID.
func (u CombUUID) Bytes() []byte {
	return u[:]
}

// Version returns algorithm version used to generate UUID.
func (u CombUUID) Version() byte {
	return u[0] >> 4
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u CombUUID) String() string {
	return UUID(u).String()
}

// SetVersion sets version bits.
func (u *CombUUID) SetVersion(v byte) {
	u[0] = (u[0] & 0x0f) | (v << 4)
}

// Returns created time in this CombUUID
func (u *CombUUID) Time() time.Time {
	t := int64(binary.BigEndian.Uint32(u[4:8]))
	t |= int64(binary.BigEndian.Uint16(u[2:4])) << 32
	t |= int64(binary.BigEndian.Uint16(u[0:4])&0xfff) << 48

	return time.Unix(0, (t-epochStart)*100)
}
