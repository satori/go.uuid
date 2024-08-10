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
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
	"testing/iotest"
	"time"
)

type faultyReader struct {
	callsNum   int
	readToFail int // Read call number to fail
}

func (r *faultyReader) Read(dest []byte) (int, error) {
	r.callsNum++
	if (r.callsNum - 1) == r.readToFail {
		return 0, fmt.Errorf("io: reader is faulty")
	}
	return rand.Read(dest)
}

func TestNewV1(t *testing.T) {
	u1, err := NewV1()
	require.NoError(t, err)
	assert.Equal(t, V1, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())

	u2, err := NewV1()
	require.NoError(t, err)
	assert.NotEqual(t, u1, u2)
}

func TestNewV1EpochStale(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc: func() time.Time {
			return time.Unix(0, 0)
		},
		hwAddrFunc: defaultHWAddrFunc,
		rand:       rand.Reader,
	}
	u1, err := g.NewV1()
	require.NoError(t, err)

	u2, err := g.NewV1()
	require.NoError(t, err)
	assert.NotEqual(t, u1, u2)
}

func TestNewV1FaultyRand(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       &faultyReader{},
	}
	u1, err := g.NewV1()
	require.Error(t, err)
	assert.Equal(t, Nil, u1)
}

func TestNewV1MissingNetworkInterfaces(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc: time.Now,
		hwAddrFunc: func() (net.HardwareAddr, error) {
			return nil, fmt.Errorf("uuid: no hw address found")
		},
		rand: rand.Reader,
	}
	u1, err := g.NewV1()
	require.NoError(t, err)
	assert.NotEqual(t, Nil, u1)
}

func TestNewV1MissingNetInterfacesAndFaultyRand(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc: time.Now,
		hwAddrFunc: func() (net.HardwareAddr, error) {
			return nil, fmt.Errorf("uuid: no hw address found")
		},
		rand: &faultyReader{
			readToFail: 1,
		},
	}
	u1, err := g.NewV1()
	require.Error(t, err)
	assert.Equal(t, Nil, u1)
}

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV1()
	}
}

func TestNewV2(t *testing.T) {
	u1, err := NewV2(DomainPerson)
	require.NoError(t, err)
	assert.Equal(t, V2, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())

	u2, err := NewV2(DomainGroup)
	require.NoError(t, err)
	assert.Equal(t, V2, u2.Version())
	assert.Equal(t, VariantRFC4122, u2.Variant())

	u3, err := NewV2(DomainOrg)
	require.NoError(t, err)
	assert.Equal(t, V2, u3.Version())
	assert.Equal(t, VariantRFC4122, u3.Variant())
}

func TestNewV2FaultyRand(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       &faultyReader{},
	}
	u1, err := g.NewV2(DomainPerson)
	require.Error(t, err)
	assert.Equal(t, Nil, u1)
}

func BenchmarkNewV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV2(DomainPerson)
	}
}

func TestNewV3(t *testing.T) {
	u1 := NewV3(NamespaceDNS, "www.example.com")
	assert.Equal(t, V3, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())
	assert.Equal(t, "5df41881-3aed-3515-88a7-2f4a814cf09e", u1.String())

	u2 := NewV3(NamespaceDNS, "example.com")
	assert.NotEqual(t, u1, u2)

	u3 := NewV3(NamespaceDNS, "example.com")
	assert.Equal(t, u2, u3)

	u4 := NewV3(NamespaceURL, "example.com")
	assert.NotEqual(t, u3, u4)
}

func BenchmarkNewV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV3(NamespaceDNS, "www.example.com")
	}
}

func TestNewV4(t *testing.T) {
	u1, err := NewV4()
	require.NoError(t, err)
	assert.Equal(t, V4, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())

	u2, err := NewV4()
	require.NoError(t, err)
	assert.NotEqual(t, u1, u2)
}

func TestNewV4FaultyRand(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       &faultyReader{},
	}
	u1, err := g.NewV4()
	require.Error(t, err)
	assert.Equal(t, Nil, u1)
}

func TestNewV4PartialRead(t *testing.T) {
	g := &rfc4122Generator{
		epochFunc:  time.Now,
		hwAddrFunc: defaultHWAddrFunc,
		rand:       iotest.OneByteReader(rand.Reader),
	}
	u1, err := g.NewV4()
	require.NoError(t, err)

	zeros := bytes.Count(u1.Bytes(), []byte{0})
	mostlyZeros := zeros >= 10

	assert.False(t, mostlyZeros, "Generated UUID contains mostly zeros")
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV4()
	}
}

func TestNewV5(t *testing.T) {
	u1 := NewV5(NamespaceDNS, "www.example.com")
	assert.Equal(t, V5, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())
	assert.Equal(t, "2ed6657d-e927-568b-95e1-2665a8aea6a2", u1.String())

	u2 := NewV5(NamespaceDNS, "example.com")
	assert.NotEqual(t, u1, u2)

	u3 := NewV5(NamespaceDNS, "example.com")
	assert.Equal(t, u2, u3)

	u4 := NewV5(NamespaceURL, "example.com")
	assert.NotEqual(t, u3, u4)
}

func BenchmarkNewV5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV5(NamespaceDNS, "www.example.com")
	}
}

func TestNewV6(t *testing.T) {
	u1, err := NewV6()
	require.NoError(t, err)
	assert.Equal(t, V6, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())

	u2, err := NewV6()
	require.NoError(t, err)
	assert.NotEqual(t, u1, u2)
	assert.True(t, bytes.Compare(u1[:6], u2[:6]) < 0 || bytes.Equal(u1[:6], u2[:6]))
}

func BenchmarkNewV6(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV6()
	}
}

func TestNewV7(t *testing.T) {
	u1, err := NewV7()
	require.NoError(t, err)
	assert.Equal(t, V7, u1.Version())
	assert.Equal(t, VariantRFC4122, u1.Variant())

	u2, err := NewV7()
	require.NoError(t, err)
	assert.NotEqual(t, u1, u2)

	assert.True(t, bytes.Compare(u1[:6], u2[:6]) < 0 || bytes.Equal(u1[:6], u2[:6]))
	assert.NotEqual(t, u1[6:], u2[6:])
}

func BenchmarkNewV7(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewV7()
	}
}
