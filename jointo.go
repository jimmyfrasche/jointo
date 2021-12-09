// Package jointo provides functions for joining strings and byte slices
// with a separator while writing the output directly to an io.Writer.
//
// If the writer has a Grow(int) method, such as strings.Builder and bytes.Buffer,
// it's used to preallocate all necessary space.
package jointo

import (
	"io"
)

func cvt(n int, err error) (int64, error) {
	return int64(n), err
}

// String writes elems joined by sep to w.
//
// It is equivalent to
//	io.WriteString(w, strings.Join(elems, sep))
func String(w io.Writer, elems []string, sep string) (written int64, err error) {
	switch len(elems) {
	case 0:
		return 0, nil
	case 1:
		return cvt(io.WriteString(w, elems[0]))
	}

	// If we can pre-allocate space in w, do so.
	if grower, ok := w.(interface{ Grow(int) }); ok {
		n := len(sep) * (len(elems) - 1)
		for i := 0; i < len(elems); i++ {
			n += len(elems[i])
		}
		grower.Grow(n)

	}

	// it's possible that the compiler is intelligent enough to allow this
	// to be written with io.WriteString, inline it and output equivalent code
	// but that would require additional testing to verify.
	if sw, ok := w.(io.StringWriter); ok {
		n, err := sw.WriteString(elems[0])
		written += int64(n)
		if err != nil {
			return written, err
		}
		for _, elem := range elems[1:] {
			n, err = sw.WriteString(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}

			n, err = sw.WriteString(elem)
			written += int64(n)
			if err != nil {
				return written, err
			}
		}
	} else {
		n, err := w.Write([]byte(elems[0]))
		written += int64(n)
		if err != nil {
			return written, err
		}
		sep := []byte(sep)
		for _, elem := range elems[1:] {
			n, err = w.Write(sep)
			written += int64(n)
			if err != nil {
				return written, err
			}

			n, err = w.Write([]byte(elem))
			written += int64(n)
			if err != nil {
				return written, err
			}
		}
	}

	return written, nil
}

// Bytes writes elems joined by sep to w.
//
// It is equivalent to
//	w.Write(w, bytes.Join(elems, sep))
func Bytes(w io.Writer, elems [][]byte, sep []byte) (written int64, err error) {
	switch len(elems) {
	case 0:
		return 0, nil
	case 1:
		return cvt(w.Write(elems[0]))
	}

	// If we can pre-allocate space in w, do so.
	if grower, ok := w.(interface{ Grow(int) }); ok {
		n := len(sep) * (len(elems) - 1)
		for i := 0; i < len(elems); i++ {
			n += len(elems[i])
		}
		grower.Grow(n)

	}

	n, err := w.Write([]byte(elems[0]))
	written += int64(n)
	if err != nil {
		return written, err
	}
	for _, elem := range elems[1:] {
		n, err = w.Write(sep)
		written += int64(n)
		if err != nil {
			return written, err
		}
		n, err = w.Write([]byte(elem))
		written += int64(n)
		if err != nil {
			return written, err
		}
	}

	return written, nil
}
