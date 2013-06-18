/*
The uuid package provides implementation of Universally Unique Identifier (UUID) structure
with functions for generating versions 3, 4 and 5 as specified in RFC 4122
*/
package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"hash"
)

// UUID layout variants.
const (
	VariantNCS = iota
	VariantRFC4122
	VariantMicrosoft
	VariantFuture
)

// UUID representation compliant with specification
// described in RFC 4122.
type UUID [16]byte

var (
	NamespaceDNS  = UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceURL  = UUID{0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceOID  = UUID{0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
	NamespaceX500 = UUID{0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}
)

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

// Returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns UUID, name string) (u *UUID, err error) {
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
func NewV5(ns UUID, name string) (u *UUID, err error) {
	u, err = newFromHash(sha1.New(), ns, name)
	if err != nil {
		return
	}
	u.setVersion(5)
	u.setVariant()
	return
}

// Returns UUID based on hashing of namespace UUID and name
func newFromHash(h hash.Hash, ns UUID, name string) (u *UUID, err error) {
	u = new(UUID)
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))
	return
}
