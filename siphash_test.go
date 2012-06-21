package siphash

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var zeroKey = make([]byte, 16)

var golden = []struct {
	k []byte
	m []byte
	r uint64
}{
	{
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
		[]byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e},
		0xa129ca6149be45e5,
	},
	{
		zeroKey,
		[]byte("Hello world"),
		0xc9e8a3021f3822d9,
	},
	{
		zeroKey,
		[]byte{}, // zero-length message
		0x1e924b9d737700d7,
	},
	{
		zeroKey,
		[]byte("12345678123"),
		0xf95d77ccdb0649f,
	},
	{
		zeroKey,
		make([]byte, 8),
		0xe849e8bb6ffe2567,
	},
	{
		zeroKey,
		make([]byte, 1535),
		0xe74d1c0ab64b2afa,
	},
}

func TestSum64(t *testing.T) {
	for i, v := range golden {
		h := New(v.k)
		h.Write(v.m)
		if sum := h.Sum64(); sum != v.r {
			t.Errorf(`%d: expected "%x", got "%x"`, i, v.r, sum)
		}
	}
}

func TestSum(t *testing.T) {
	var r [8]byte
	for i, v := range golden {
		binary.LittleEndian.PutUint64(r[:], v.r)
		h := New(v.k)
		h.Write(v.m)
		if sum := h.Sum(nil); !bytes.Equal(sum, r[:]) {
			t.Errorf(`%d: expected "%x", got "%x"`, i, r, sum)
		}
	}
}

var key = zeroKey
var bench = New(key)
var buf = make([]byte, 8<<10)

func BenchmarkHash8(b *testing.B) {
	b.SetBytes(8)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:8])
		bench.Sum64()
	}
}

func BenchmarkHash16(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:16])
		bench.Sum64()
	}
}

func BenchmarkHash40(b *testing.B) {
	b.SetBytes(24)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:16])
		bench.Sum64()
	}
}

func BenchmarkHash64(b *testing.B) {
	b.SetBytes(64)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:64])
		bench.Sum64()
	}
}

func BenchmarkHash128(b *testing.B) {
	b.SetBytes(128)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:64])
		bench.Sum64()
	}
}

func BenchmarkHash1K(b *testing.B) {
	b.SetBytes(1024)
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf[:1024])
		bench.Sum64()
	}
}

func BenchmarkHash8K(b *testing.B) {
	b.SetBytes(int64(len(buf)))
	for i := 0; i < b.N; i++ {
		bench.Reset()
		bench.Write(buf)
		bench.Sum64()
	}
}
