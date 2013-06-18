/*
The uuid package provides implementation of Universally Unique Identifier (UUID) structure
with functions for generating versions 1, 3, 4 and 5 as specified in RFC 4122
*/
package uuid

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"net"
	"time"
)

// UUID layout variants.
const (
	VariantNCS = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// Difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
const epochStart = 122192928000000000

// Clock sequence storage
var clockSequence uint16

// UUID representation compliant with specification
// described in RFC 4122.
type UUID [16]byte

// Predefined namespace UUIDs.
var (
	NamespaceDNS  = &UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceURL  = &UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceOID  = &UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceX500 = &UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
)

// Returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 *UUID, u2 *UUID) bool {
	return bytes.Equal(u1[:], u2[:])
}

// Returns algorithm version used to generate UUID.
// RFC 4122 describes version 1, 3, 4 and 5.
func (u *UUID) Version() uint {
	return uint(u[6] >> 4)
}

// Returns UUID layout variant.
func (u *UUID) Variant() uint {
	switch {
	case (u[8] & 0x80) == 0x00:
		return VariantNCS
	case (u[8]&0xc0)|0x80 == 0x80:
		return VariantRFC4122
	case (u[8]&0xe0)|0xc0 == 0xc0:
		return VariantMicrosoft
	}
	return VariantFuture
}

// Returns canonical string representation of UUID:
// xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx.
func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

// Sets version bits.
func (u *UUID) setVersion(v byte) {
	u[6] = (u[6] & 0x0f) | (v << 4)
}

// Sets variant bits as described in RFC 4122.
func (u *UUID) setVariant() {
	u[8] = (u[8] & 0xbf) | 0x80
}

// Returns UUID based on current timestamp and MAC address.
func NewV1() (u *UUID, err error) {
	u = new(UUID)

	time_now := epochStart + uint64(time.Now().UnixNano()/100)

	binary.BigEndian.PutUint32(u[0:], uint32(time_now))
	binary.BigEndian.PutUint16(u[4:], uint16(time_now>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(time_now>>48))

	if clockSequence == 0 {
		buf := make([]byte, 2)
		_, err = rand.Read(buf)
		clockSequence = binary.BigEndian.Uint16(buf)
	}

	binary.BigEndian.PutUint16(u[8:], clockSequence)

	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		if len(iface.HardwareAddr) >= 6 {
			copy(u[10:], iface.HardwareAddr)
			break
		}
	}

	u.setVersion(1)
	u.setVariant()
	return
}

// Returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns *UUID, name string) (u *UUID, err error) {
	u, err = newFromHash(md5.New(), ns, name)
	if err != nil {
		return
	}
	u.setVersion(3)
	u.setVariant()
	return
}

// Returns random generated UUID.
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

// Returns UUID based on SHA-1 hash of namespace UUID and name.
func NewV5(ns *UUID, name string) (u *UUID, err error) {
	u, err = newFromHash(sha1.New(), ns, name)
	if err != nil {
		return
	}
	u.setVersion(5)
	u.setVariant()
	return
}

// Returns UUID based on hashing of namespace UUID and name.
func newFromHash(h hash.Hash, ns *UUID, name string) (u *UUID, err error) {
	u = new(UUID)
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))
	return
}
