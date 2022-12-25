package transport

import (
	"bytes"
	"fmt"
	"io"
)

func DeserializeSlice(r io.Reader) ([]byte, error) {
	length := make([]byte, 1)
	n, err := r.Read(length)
	if err != nil {
		return nil, fmt.Errorf("cant read field len: %v", err)
	}
	if n != 1 {
		return nil, fmt.Errorf("field with unexpected length %d", n)
	}

	if length[0] == 0 {
		return nil, nil
	}

	byteSlice := make([]byte, length[0])
	n, err = r.Read(byteSlice)
	if err != nil {
		return nil, fmt.Errorf("cant read byte slice: %v", err)
	}
	if byte(n) != length[0] {
		return nil, fmt.Errorf("byte slice with unexpected length %d", n)
	}

	return byteSlice, nil
}

func SerializeSlice(slice []byte) []byte {
	if len(slice) > 255 {
		return nil
	}

	return bytes.Join(
		[][]byte{
			{byte(len(slice))},
			slice,
		},
		nil)
}

func WriteSlice(w io.Writer, slice []byte) error {
	n, err := w.Write(slice)
	if err != nil {
		return fmt.Errorf("cant write: %w", err)
	}
	if n != len(slice) {
		return fmt.Errorf("write %d of %d ", n, len(slice))
	}

	return nil
}
