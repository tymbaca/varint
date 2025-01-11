package varint

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_encodeZigZag64(t *testing.T) {
	cases := []struct {
		arg    int64
		expect uint64
	}{
		{arg: 0, expect: 0}, // source: https://protobuf.dev/programming-guides/encoding/#signed-ints
		{arg: -1, expect: 1},
		{arg: 1, expect: 2},
		{arg: -2, expect: 3},
		{arg: 2, expect: 4},
		{arg: 0x7fffffff, expect: 0xfffffffe},
		{arg: -0x80000000, expect: 0xffffffff},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			require.Equal(t, tt.expect, encodeZigZag(tt.arg))
		})
	}
}

func Test_decodeZigZag64(t *testing.T) {
	cases := []struct {
		arg    uint64
		expect int64
	}{
		{expect: 0, arg: 0}, // source: https://protobuf.dev/programming-guides/encoding/#signed-ints
		{expect: -1, arg: 1},
		{expect: 1, arg: 2},
		{expect: -2, arg: 3},
		{expect: 2, arg: 4},
		{expect: 0x7fffffff, arg: 0xfffffffe},
		{expect: -0x80000000, arg: 0xffffffff},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			require.Equal(t, tt.expect, decodeZigZag(tt.arg))
		})
	}
}
