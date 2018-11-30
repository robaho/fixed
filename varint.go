package fixed

import "io"

// The binary encoding package does not offer 'write' methods so the allocation costs in WriteTo are high
// so we duplicate the code here and implement them

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.

func putUvarint(w io.ByteWriter, x uint64) error {
	i := 0
	for x >= 0x80 {
		err := w.WriteByte(byte(x) | 0x80)
		if err != nil {
			return err
		}
		x >>= 7
		i++
	}
	return w.WriteByte(byte(x))
}

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.

func putVarint(w io.ByteWriter, x int64) error {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return putUvarint(w, ux)
}
