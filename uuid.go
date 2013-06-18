package uuid

import (
	"fmt"
	"crypto/rand"
)

type UUID [16]byte

// Returns algorithm version used to generate UUID
// RFC 4122 describes version 1, 3, 4 and 5.
func (u *UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// Set version bits.
func (u *UUID) setVersion(v uint) {
	u[6] = (u[6] & 0x0f) | (byte(v) << 4)
}

// Set variant bits.
func (u *UUID) setVariant() {
	u[8] = (u[8] & 0x3f) | 0x80
}

// Returns random UUID.
func NewV4() (u *UUID, err error) {
	u = new(UUID)
	_, err = rand.Read(u[:])
	if err != nil {
		return
	}
	u.setVersion(4)
	u.setVariant()
	return
}
