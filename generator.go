package uuid

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"sync"
	"hash"
	"net"
	"io"
	"os"
	"time"
	"crypto/rand"
)


var (
	global = NewGenerator(rand.Reader)

	epochFunc = unixTimeFunc
	posixUID  = uint32(os.Getuid())
	posixGID  = uint32(os.Getgid())
)

// Generator provides interface for generating UUIDs
type Generator interface {
	NewV1() UUID
	NewV2(domain byte) UUID
	NewV3(ns UUID, name string) UUID
	NewV4() UUID
	NewV5(ns UUID, name string) UUID
}

// NewGenerator creates new generator given the source of randomness (use "crypto/rand".Reader)
func NewGenerator(rand io.Reader) Generator {
	return &generator{
		rand: rand,
	}
}

// NewV1 returns UUID based on current timestamp and MAC address.
func NewV1() UUID {
	return global.NewV1()
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func NewV2(domain byte) UUID {
	return global.NewV2(domain)
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func NewV3(ns UUID, name string) UUID {
	return global.NewV3(ns, name)
}

// NewV4 returns random generated UUID.
func NewV4() UUID {
	return global.NewV4()
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func NewV5(ns UUID, name string) UUID {
	return global.NewV5(ns, name)
}

type generator struct {
	storageOnce sync.Once
	storageMutex sync.Mutex

	rand io.Reader

	lastTime      uint64
	clockSequence uint16
	hardwareAddr  [6]byte
}

// NewV1 returns UUID based on current timestamp and MAC address.
func (g *generator) NewV1() UUID {
	u := UUID{}

	timeNow, clockSeq, hardwareAddr := g.getStorage()

	binary.BigEndian.PutUint32(u[0:], uint32(timeNow))
	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)

	copy(u[10:], hardwareAddr)

	u.SetVersion(1)
	u.SetVariant()

	return u
}

// NewV2 returns DCE Security UUID based on POSIX UID/GID.
func (g *generator) NewV2(domain byte) UUID {
	u := UUID{}

	timeNow, clockSeq, hardwareAddr := g.getStorage()

	switch domain {
	case DomainPerson:
		binary.BigEndian.PutUint32(u[0:], posixUID)
	case DomainGroup:
		binary.BigEndian.PutUint32(u[0:], posixGID)
	}

	binary.BigEndian.PutUint16(u[4:], uint16(timeNow>>32))
	binary.BigEndian.PutUint16(u[6:], uint16(timeNow>>48))
	binary.BigEndian.PutUint16(u[8:], clockSeq)
	u[9] = domain

	copy(u[10:], hardwareAddr)

	u.SetVersion(2)
	u.SetVariant()

	return u
}

// NewV3 returns UUID based on MD5 hash of namespace UUID and name.
func (g *generator) NewV3(ns UUID, name string) UUID {
	u := newFromHash(md5.New(), ns, name)
	u.SetVersion(3)
	u.SetVariant()

	return u
}

// NewV4 returns random generated UUID.
func (g *generator) NewV4() UUID {
	u := UUID{}
	g.safeRandom(u[:])
	u.SetVersion(4)
	u.SetVariant()

	return u
}

// NewV5 returns UUID based on SHA-1 hash of namespace UUID and name.
func (g *generator) NewV5(ns UUID, name string) UUID {
	u := newFromHash(sha1.New(), ns, name)
	u.SetVersion(5)
	u.SetVariant()

	return u
}

func (g *generator) initStorage() {
	g.initClockSequence()
	g.initHardwareAddr()
}

func (g *generator) initClockSequence() {
	buf := make([]byte, 2)
	g.safeRandom(buf)
	g.clockSequence = binary.BigEndian.Uint16(buf)
}

func (g *generator) initHardwareAddr() {
	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if len(iface.HardwareAddr) >= 6 {
				copy(g.hardwareAddr[:], iface.HardwareAddr)
				return
			}
		}
	}

	// Initialize hardwareAddr randomly in case
	// of real network interfaces absence
	g.safeRandom(g.hardwareAddr[:])

	// Set multicast bit as recommended in RFC 4122
	g.hardwareAddr[0] |= 0x01
}

func (g *generator) safeRandom(dest []byte) {
	if _, err := g.rand.Read(dest); err != nil {
		panic(err)
	}
}

// Returns UUID v1/v2 storage state.
// Returns epoch timestamp, clock sequence, and hardware address.
func (g *generator) getStorage() (uint64, uint16, []byte) {
	g.storageOnce.Do(g.initStorage)

	g.storageMutex.Lock()
	defer g.storageMutex.Unlock()

	timeNow := epochFunc()
	// Clock changed backwards since last UUID generation.
	// Should increase clock sequence.
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	return timeNow, g.clockSequence, g.hardwareAddr[:]
}

// Returns difference in 100-nanosecond intervals between
// UUID epoch (October 15, 1582) and current time.
// This is default epoch calculation function.
func unixTimeFunc() uint64 {
	return epochStart + uint64(time.Now().UnixNano()/100)
}

// Returns UUID based on hashing of namespace UUID and name.
func newFromHash(h hash.Hash, ns UUID, name string) UUID {
	u := UUID{}
	h.Write(ns[:])
	h.Write([]byte(name))
	copy(u[:], h.Sum(nil))

	return u
}
