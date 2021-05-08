package uuid

import (
	"bytes"
	"crypto/rand"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func TestCombUUID(t *testing.T) { TestingT(t) }

type testSuiteComb struct{}

var _ = Suite(&testSuiteComb{})

func (s *testSuiteComb) TestBytes(c *C) {
	u := UUID{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	bytes1 := []byte{0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8}

	c.Assert(bytes.Equal(u.Bytes(), bytes1), Equals, true)
}

func (s *testSuiteComb) TestString(c *C) {
	c.Assert(NamespaceDNS.String(), Equals, "6ba7b810-9dad-11d1-80b4-00c04fd430c8")
}

func (s *testSuiteComb) TestEqual(c *C) {
	c.Assert(CombUUID(NamespaceDNS).Equal(CombUUID(NamespaceDNS)), Equals, true)
	c.Assert(CombUUID(NamespaceDNS).Equal(CombUUID(NamespaceURL)), Equals, false)
}

func (s *testSuiteComb) TestVersion(c *C) {
	u := UUID{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	c.Assert(u.Version(), Equals, V1)
}

func (s *testSuiteComb) TestSetVersion(c *C) {
	u := UUID{}
	u.SetVersion(4)
	c.Assert(u.Version(), Equals, V4)
}

func (s *testSuiteComb) TestTime(c *C) {
	now := time.Unix(0, 0)

	g := &rfc4122AndCombGenerator{
		epochFunc: func() time.Time {
			return now
		},
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}

	u1, err := g.NewCombV4()
	c.Assert(err, IsNil)

	if err == nil {
		ut := u1.Time()
		c.Assert(ut.Unix(), Equals, now.Unix())
	}

	u2, err := g.NewCombV1()
	c.Assert(err, IsNil)

	if err == nil {
		ut := u2.Time()
		c.Assert(ut.Unix(), Equals, now.Unix())
	}
}
