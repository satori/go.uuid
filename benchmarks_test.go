package uuid

import (
	"testing"
)

func BenchmarkNewV1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV1()
	}
}

func BenchmarkNewV2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV2(DomainPerson)
	}
}

func BenchmarkNewV3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV3(NamespaceDNS, "www.example.com")
	}
}

func BenchmarkNewV4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV4()
	}
}

func BenchmarkNewV5(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewV5(NamespaceDNS, "www.example.com")
	}
}
