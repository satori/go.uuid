// Copyright (C) 2013-2018 by Maxim Bublis <b@codemonkey.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package uuid

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// Difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and Unix epoch (January 1, 1970).
const epochStart = 122192928000000000

type epochFunc func() time.Time
type hwAddrFunc func() (net.HardwareAddr, error)

var (
	global = newRfc4122AndCombGenerator()

	posixUID = uint32(os.Getuid())
	posixGID = uint32(os.Getgid())
)

// NewV1 returns UUID based on current timestamp and MAC address.
func NewV1() (UUID, error) {
	return global.NewV1()
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func NewV2(domain byte) (UUID, error) {
	return global.NewV2(domain)
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	return global.NewV3(ns, name)
}

// NewV4 returns random generated UUID.
func NewV4() (UUID, error) {
	return global.NewV4()
}

// NewCombV1 returns UUID based on current timestamp and MAC address with ordered-time.
func NewCombV1() (UUID, error) {
	return global.NewCombV1()
}

// NewCombV4 returns UUID based on current timestamp and random generated UUID with ordered-time.
func NewCombV4() (UUID, error) {
	return global.NewCombV4()
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	return global.NewV5(ns, name)
}

// Generator provides interface for generating UUIDs.
type Generator interface {
	// RFC 4122
	NewV1() (UUID, error)
	NewV2(domain byte) (UUID, error)
	NewV3(ns UUID, name string) UUID
	NewV4() (UUID, error)
	NewV5(ns UUID, name string) UUID

	// COMB
	NewCombV1() (UUID, error)
	NewCombV4() (UUID, error)
}

// Default generator implementation.
type rfc4122AndCombGenerator struct {
	clockSequenceOnce sync.Once
	hardwareAddrOnce  sync.Once
	storageMutex      sync.Mutex

	rand io.Reader

	epochFunc     epochFunc
	hwAddrFunc    hwAddrFunc
	lastTime      uint64
	clockSequence uint16
	hardwareAddr  [6]byte
}

func newRfc4122AndCombGenerator() Generator {
	return &rfc4122AndCombGenerator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
}

// NewV1 returns UUID based on current timestamp and MAC address.
func (g *rfc4122AndCombGenerator) NewV1() (UUID, error) {
	u := UUID{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return Nil, err
	}
	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	hardwareAddr, err := g.getHardwareAddr()
	if err != nil {
		return Nil, err
	}
	copy(u[10:], hardwareAddr)

	u.SetVersion(V1)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func (g *rfc4122AndCombGenerator) NewV2(domain byte) (UUID, error) {
	u, err := g.NewV1()
	if err != nil {
		return Nil, err
	}

	switch domain {
	case DomainPerson:
		binary.BigEndian.PutUint32(u[:], posixUID)
	case DomainGroup:
		binary.BigEndian.PutUint32(u[:], posixGID)
	}

	u[9] = domain

	u.SetVersion(V2)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func (g *rfc4122AndCombGenerator) NewV3(ns UUID, name string) UUID {
	u := newFromHash(md5.New(), ns, name)
	u.SetVersion(V3)
	u.SetVariant(VariantRFC4122)

	return u
}

// NewV4 returns random generated UUID.
func (g *rfc4122AndCombGenerator) NewV4() (UUID, error) {
	u := UUID{}
	if _, err := io.ReadFull(g.rand, u[:]); err != nil {
		return Nil, err
	}
	u.SetVersion(V4)
	u.SetVariant(VariantRFC4122)

	return u, nil
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func (g *rfc4122AndCombGenerator) NewV5(ns UUID, name string) UUID {
	u := newFromHash(sha1.New(), ns, name)
	u.SetVersion(V5)
	u.SetVariant(VariantRFC4122)

	return u
}

// NewCombV1 returns nonstandard UUID based on V1 RFC4122 with timestamp bytes in the front
func (g *rfc4122AndCombGenerator) NewCombV1() (UUID, error) {
	u := UUID{}

	timeNow, clockSeq, err := g.getClockSequence()
	if err != nil {
		return Nil, err
	}

	binary.BigEndian.PutUint16(u[0:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[2:], uint16(timeNow>>32))
	binary.BigEndian.PutUint32(u[4:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	hardwareAddr, err := g.getHardwareAddr()
	if err != nil {
		return Nil, err
	}
	copy(u[10:], hardwareAddr)

	u.SetVersion(V1)
	u.SetVariant(VariantFuture)

	return u, nil
}

// NewCombV4 returns nonstandard UUID based on V4 RFC4122
// with timestamp bytes in the front (not fully random)
func (g *rfc4122AndCombGenerator) NewCombV4() (UUID, error) {
	u := UUID{}

	timeNow := g.getEpoch()

	binary.BigEndian.PutUint16(u[0:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[2:], uint16(timeNow>>32))
	binary.BigEndian.PutUint32(u[4:], uint32(timeNow))

	if _, err := io.ReadFull(g.rand, u[8:]); err != nil {
		return Nil, err
	}

	u.SetVersion(V4)
	u.SetVariant(VariantFuture)

	return u, nil
}

// Returns epoch and clock sequence.
func (g *rfc4122AndCombGenerator) getClockSequence() (uint64, uint16, error) {
	var err error
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		if _, err = io.ReadFull(g.rand, buf); err != nil {
			return
		}
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})
	if err != nil {
		return 0, 0, err
	}

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := g.getEpoch()
	// Clock didn't change since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, nil
}

// Returns hardware address.
func (g *rfc4122AndCombGenerator) getHardwareAddr() ([]byte, error) {
	var err error
	g.hardwareAddrOnce.Do(func() {
		if hwAddr, err := g.hwAddrFunc(); err == nil {
			copy(g.hardwareAddr[:], hwAddr)
			return
		}

		// Initialize hardwareAddr randomly in case
		// of real network interfaces absence.
		if _, err = io.ReadFull(g.rand, g.hardwareAddr[:]); err != nil {
			return
		}
		// Set multicast bit as recommended by RFC 4122
		g.hardwareAddr[0] |= 0x01
	})
	if err != nil {
		return []byte{}, err
	}
	return g.hardwareAddr[:], nil
}

// Returns difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and current time.
func (g *rfc4122AndCombGenerator) getEpoch() uint64 {
	return epochStart + uint64(g.epochFunc().UnixNano()/100)
}

// Returns UUID based on hashing of namespace UUID and name.
func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}

// Returns hardware address.
func defaultHWAddrFunc() (net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return []byte{}, err
	}
	for _, iface := range ifaces {
		if len(iface.HardwareAddr) >= 6 {
			return iface.HardwareAddr, nil
		}
	}
	return []byte{}, fmt.Errorf("uuid: no HW address found")
}
