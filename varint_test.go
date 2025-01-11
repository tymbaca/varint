package varint

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadUint64(t *testing.T) {
	cases := []struct {
		expect uint64
		arg    []byte
	}{
		{expect: 0, arg: []byte{0b00000000}},
		{expect: 1, arg: []byte{0b00000001}},
		{expect: 127, arg: []byte{0b01111111}},
		{expect: 128, arg: []byte{0b10000000, 0b00000001}},
		{expect: 128, arg: []byte{0b10000000, 0b10000001, 0b00000000}},
		{expect: 128, arg: []byte{0b10000000, 0b10000001, 0b10000000, 0b00000000}},
		{expect: 129, arg: []byte{0b10000001, 0b00000001}},
		{expect: 256, arg: []byte{0b10000000, 0b00000010}},
		{expect: 1024, arg: []byte{0b10000000, 0b00001000}},
		// TODO: try too much bytes
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			actual, n, err := ReadUint64(bytes.NewBuffer(tt.arg))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
			require.Equal(t, len(tt.arg), n)
		})
	}
}

func TestReadUint32(t *testing.T) {
	cases := []struct {
		expect uint32
		arg    []byte
	}{
		{expect: 0, arg: []byte{0b00000000}},
		{expect: 1, arg: []byte{0b00000001}},
		{expect: 127, arg: []byte{0b01111111}},
		{expect: 128, arg: []byte{0b10000000, 0b00000001}},
		{expect: 128, arg: []byte{0b10000000, 0b10000001, 0b00000000}},
		{expect: 128, arg: []byte{0b10000000, 0b10000001, 0b10000000, 0b00000000}},
		{expect: 129, arg: []byte{0b10000001, 0b00000001}},
		{expect: 256, arg: []byte{0b10000000, 0b00000010}},
		{expect: 1024, arg: []byte{0b10000000, 0b00001000}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			actual, n, err := ReadUint32(bytes.NewBuffer(tt.arg))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
			require.Equal(t, len(tt.arg), n)
		})
	}
}

func TestReadInt64(t *testing.T) {
	cases := []struct {
		expect int64
		arg    []byte
	}{
		{expect: 0, arg: []byte{0b00000000}},
		{expect: -1, arg: []byte{0b00000001}},
		{expect: 1, arg: []byte{0b00000010}},
		{expect: -2, arg: []byte{0b00000011}},

		// 0x7fffffff (2147483647) is decoded from 0xfffffffe (0b11111111111111111111111111111110)
		//
		// 0001111 1111111 1111111 1111111 1111110
		// 0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{expect: 2147483647, arg: []byte{0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},

		// -0x80000000 (-2147483648) is decoded from 0xffffffff (0b11111111111111111111111111111111)
		//
		// 0001111 1111111 1111111 1111111 1111111
		// 0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{expect: -0x80000000, arg: []byte{0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			actual, n, err := ReadInt64(bytes.NewBuffer(tt.arg))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
			require.Equal(t, len(tt.arg), n)
		})
	}
}

func TestReadInt32(t *testing.T) {
	cases := []struct {
		expect int32
		arg    []byte
	}{
		{expect: 0, arg: []byte{0b00000000}},
		{expect: -1, arg: []byte{0b00000001}},
		{expect: 1, arg: []byte{0b00000010}},
		{expect: -2, arg: []byte{0b00000011}},

		// 0x7fffffff (2147483647) is decoded from 0xfffffffe (0b11111111111111111111111111111110)
		//
		// 0001111 1111111 1111111 1111111 1111110
		// 0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{expect: 2147483647, arg: []byte{0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},

		// -0x80000000 (-2147483648) is decoded from 0xffffffff (0b11111111111111111111111111111111)
		//
		// 0001111 1111111 1111111 1111111 1111111
		// 0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{expect: -0x80000000, arg: []byte{0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			actual, n, err := ReadInt32(bytes.NewBuffer(tt.arg))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
			require.Equal(t, len(tt.arg), n)
		})
	}
}

func TestWriteUint64(t *testing.T) {
	cases := []struct {
		arg    uint64
		expect []byte
	}{
		{arg: 0, expect: []byte{0b00000000}},
		{arg: 1, expect: []byte{0b00000001}},
		{arg: 127, expect: []byte{0b01111111}},
		{arg: 128, expect: []byte{0b10000000, 0b00000001}},
		{arg: 129, expect: []byte{0b10000001, 0b00000001}},
		{arg: 256, expect: []byte{0b10000000, 0b00000010}},
		{arg: 1024, expect: []byte{0b10000000, 0b00001000}},
		// TODO: try too much bytes
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			n, err := WriteUint64(buf, tt.arg)
			require.NoError(t, err)
			require.Equal(t, tt.expect, buf.Bytes())
			require.Equal(t, len(tt.expect), n)
		})
	}
}

func TestWriteUint32(t *testing.T) {
	cases := []struct {
		arg    uint32
		expect []byte
	}{
		{arg: 0, expect: []byte{0b00000000}},
		{arg: 1, expect: []byte{0b00000001}},
		{arg: 127, expect: []byte{0b01111111}},
		{arg: 128, expect: []byte{0b10000000, 0b00000001}},
		{arg: 129, expect: []byte{0b10000001, 0b00000001}},
		{arg: 256, expect: []byte{0b10000000, 0b00000010}},
		{arg: 1024, expect: []byte{0b10000000, 0b00001000}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			n, err := WriteUint32(buf, tt.arg)
			require.NoError(t, err)
			require.Equal(t, tt.expect, buf.Bytes())
			require.Equal(t, len(tt.expect), n)
		})
	}
}

func TestWriteInt64(t *testing.T) {
	cases := []struct {
		arg    int64
		expect []byte
	}{
		{arg: 0, expect: []byte{0b00000000}},
		{arg: -1, expect: []byte{0b00000001}},
		{arg: 1, expect: []byte{0b00000010}},
		{arg: -2, expect: []byte{0b00000011}},

		// 0x7fffffff (2147483647) is decoded from 0xfffffffe (0b11111111111111111111111111111110)
		//
		// 0001111 1111111 1111111 1111111 1111110
		// 0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{arg: 2147483647, expect: []byte{0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},

		// -0x80000000 (-2147483648) is decoded from 0xffffffff (0b11111111111111111111111111111111)
		//
		// 0001111 1111111 1111111 1111111 1111111
		// 0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{arg: -0x80000000, expect: []byte{0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			n, err := WriteInt64(buf, tt.arg)
			require.NoError(t, err)
			require.Equal(t, tt.expect, buf.Bytes())
			require.Equal(t, len(tt.expect), n)
		})
	}
}

func TestWriteInt32(t *testing.T) {
	cases := []struct {
		arg    int32
		expect []byte
	}{
		{arg: 0, expect: []byte{0b00000000}},
		{arg: -1, expect: []byte{0b00000001}},
		{arg: 1, expect: []byte{0b00000010}},
		{arg: -2, expect: []byte{0b00000011}},

		// 0x7fffffff (2147483647) is decoded from 0xfffffffe (0b11111111111111111111111111111110)
		//
		// 0001111 1111111 1111111 1111111 1111110
		// 0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{arg: 2147483647, expect: []byte{0b11111110, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},

		// -0x80000000 (-2147483648) is decoded from 0xffffffff (0b11111111111111111111111111111111)
		//
		// 0001111 1111111 1111111 1111111 1111111
		// 0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111,
		{arg: -0x80000000, expect: []byte{0b11111111, 0b11111111, 0b11111111, 0b11111111, 0b00001111}},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			n, err := WriteInt32(buf, tt.arg)
			require.NoError(t, err)
			require.Equal(t, tt.expect, buf.Bytes())
			require.Equal(t, len(tt.expect), n)
		})
	}
}

func Test_readBytes(t *testing.T) {
	cases := []struct {
		arg    []byte
		expect []byte
	}{
		{
			arg:    []byte{0b10000000, 0b10000001, 0b10000010, 0b00000011, 0b00000011, 0b10000011},
			expect: []byte{0b10000000, 0b10000001, 0b10000010, 0b00000011},
		},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			actual, n, err := readBytes(bytes.NewBuffer(tt.arg))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
			require.Equal(t, len(tt.expect), n)
		})
	}
}
