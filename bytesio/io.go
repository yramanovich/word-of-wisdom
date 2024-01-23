package bytesio

import (
	"bufio"
	"io"
)

// Writeln writes bytes to the underlying data stream and adds '\n' at the end.
func Writeln(w io.Writer, data []byte) error {
	_, err := w.Write(append(data, '\n'))
	return err
}

// Readln reads until the first occurrence of '\n'
// returning a slice containing the data and excluding the delimiter.
func Readln(r io.Reader) ([]byte, error) {
	data, err := bufio.NewReader(r).ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return data[:len(data)-1], nil
}
