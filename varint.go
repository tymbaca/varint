// Package varint is shit. Literally bullshit. I want zigzag LEB128 af
package varint

import (
	"bufio"
	"encoding/binary"
	"io"
)

type Int32 int32

func (i *Int32) MarshalBinary(w io.Writer, _ binary.ByteOrder) (int, error) {
	return WriteInt32(w, int32(*i))
}

func (i *Int32) UnmarshalBinary(r io.Reader, _ binary.ByteOrder) (int, error) {
	val, n, err := ReadInt32(r)
	if err != nil {
		return n, err
	}

	*i = Int32(val)
	return n, nil
}

type Int64 int64

func (i *Int64) MarshalBinary(w io.Writer, _ binary.ByteOrder) (int, error) {
	return WriteInt64(w, int64(*i))
}

func (i *Int64) UnmarshalBinary(r io.Reader, _ binary.ByteOrder) (int, error) {
	val, n, err := ReadInt64(r)
	if err != nil {
		return n, err
	}

	*i = Int64(val)
	return n, nil
}

type Uint32 uint32

func (i *Uint32) MarshalBinary(w io.Writer, _ binary.ByteOrder) (int, error) {
	return WriteUint32(w, uint32(*i))
}

func (i *Uint32) UnmarshalBinary(r io.Reader, _ binary.ByteOrder) (int, error) {
	val, n, err := ReadUint32(r)
	if err != nil {
		return n, err
	}

	*i = Uint32(val)
	return n, nil
}

type Uint64 uint64

func (i *Uint64) MarshalBinary(w io.Writer, _ binary.ByteOrder) (int, error) {
	return WriteUint64(w, uint64(*i))
}

func (i *Uint64) UnmarshalBinary(r io.Reader, _ binary.ByteOrder) (int, error) {
	val, n, err := ReadUint64(r)
	if err != nil {
		return n, err
	}

	*i = Uint64(val)
	return n, nil
}

// ReadInt32 uses protobuf zig-zag encoding
func ReadInt32(r io.Reader) (int32, int, error) {
	val, n, err := decodePayloadFrom(r)
	return int32(decodeZigZag(val)), n, err
}

// ReadInt64 uses protobuf zig-zag encoding
func ReadInt64(r io.Reader) (int64, int, error) {
	val, n, err := decodePayloadFrom(r)
	return decodeZigZag(val), n, err
}

// ReadUint32 uses basic unsigned LAB128 encoding
func ReadUint32(r io.Reader) (uint32, int, error) {
	val, n, err := decodePayloadFrom(r)
	return uint32(val), n, err
}

// ReadUint64 uses basic unsigned LAB128 encoding
func ReadUint64(r io.Reader) (uint64, int, error) {
	return decodePayloadFrom(r)
}

func WriteUint32(w io.Writer, v uint32) (int, error) {
	return encodePayloadInto(w, uint64(v))
}

func WriteUint64(w io.Writer, v uint64) (int, error) {
	return encodePayloadInto(w, v)
}

func WriteInt32(w io.Writer, v int32) (int, error) {
	bits := encodeZigZag(int64(v))
	return encodePayloadInto(w, bits)
}

func WriteInt64(w io.Writer, v int64) (int, error) {
	bits := encodeZigZag(v)
	return encodePayloadInto(w, bits)
}

func decodePayloadFrom(r io.Reader) (uint64, int, error) {
	bytes, n, err := readBytes(r)
	if err != nil {
		return 0, 0, err
	}

	return decodePayloadBits(bytes), n, nil
}

// 10101010, 01010101
//
//	0101010,  1010101
//	 |          |
//	 |    +-----+
//	 +----|------+
//	      |      |
//	  1010101,0101010
func decodePayloadBits(bytes []byte) uint64 {
	var bits uint64
	for i, b := range bytes {
		payload := uint64(dropMSB(b))
		offset := 7 * uint64(i)
		bits = bits | payload<<offset
	}

	return bits
}

func encodePayloadInto(w io.Writer, bits uint64) (int, error) {
	payload := encodePayloadBits(bits)
	return w.Write(payload)
}

// we want this behavior
// 16 bit integer:
//
//	 0010101010101010
//	   10101010101010
//	  1010101 0101010
//	      |      |
//	 +----|------+
//	 |    +-----+
//	 |          |
//	0101010,  1010101
//
// 10101010, 01010101 resulting bytes
func encodePayloadBits(bits uint64) []byte {
	// so we do this
	// 0010101010101010
	// 001010101 cut 7  payload: 0101010
	// >> 7   001010101
	//        001010101
	//        00 cut 7  payload: 1010101
	// >> 7          00
	//         00000000 (there will be zeros)
	//         0 cut 7  payload: 0000000 (this payload is junk, we can ignore it)
	// 7 + 7 + 7 = 21; 21 >= 16; break
	payloads := make([]byte, 0, 8)
	for cutBits := 0; cutBits <= 64; cutBits += 7 {
		payload := bits & 0b0111_1111 // cut out first 7 bits

		bits >>= 7 // shit original value, to prepare for next iter

		payloads = append(payloads, byte(payload))
	}

	// clean trailing zero payloads
	for i := len(payloads) - 1; i > 0; i-- { // i > 0 because we need at least 1 byte, even if it == 0
		if payloads[i] == 0 {
			payloads = payloads[:i]
		}
	}

	// mark payloads. All payloads besides last one must have 1 as the 8'th bit (MSB)
	for i := range payloads {
		// don't mark last payload
		if i == len(payloads)-1 {
			break
		}

		payloads[i] |= 0b1000_0000
	}

	return payloads
}

// readBytes reads bytes until first bit in 0 (that byte is also added to result)
func readBytes(r io.Reader) ([]byte, int, error) {
	var bytes []byte

	br := bufio.NewReader(r)
	for {
		b, err := br.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		bytes = append(bytes, b)

		if lastByte(b) {
			break
		}
	}
	return bytes, len(bytes), nil
}

func lastByte(b byte) bool {
	return b>>7 == 0
}

func dropMSB(b byte) byte {
	return b & 0b01111111
}
