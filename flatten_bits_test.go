package simdjson

import (
	"reflect"
	"testing"
)

func TestFlattenBits(t *testing.T) {

	testCases := []struct {
		bits     uint64
		expected []uint32
	}{
		{0x11,[]uint32{0x0, 0x4}},
		{0x100100100100,[]uint32{0x8, 0x14, 0x20, 0x2c}},
		{0x8101010101010101,[]uint32{0x0, 0x8, 0x10, 0x18, 0x20, 0x28, 0x30, 0x38, 0x3f}},
		{0xf000000000000000,[]uint32{0x3c, 0x3d, 0x3e, 0x3f}},
		{0xffffffffffffffff,[]uint32{
			0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf,
			0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
			0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f,
			0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x3a, 0x3b, 0x3c, 0x3d, 0x3e, 0x3f,
		}},
	}

	for i, tc := range testCases {

		base := make([]uint32, 0, 1024)

		flatten_bits(&base, uint64(64), tc.bits)

		if !reflect.DeepEqual(base, tc.expected) {
			t.Errorf("TestFlattenBits(%d): got: %v want: %v", i, base, tc.expected)
		}
	}
}
