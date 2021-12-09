package jointo

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

var table = []*struct {
	Elems []string
	Sep   string
}{
	{},
	{
		Elems: []string{"x"},
	},
	{
		Elems: []string{"x", "y", "z"},
		Sep:   "/",
	},
}

type noGrow struct {
	w *bytes.Buffer
}

func (w *noGrow) Write(v []byte) (int, error) {
	return w.w.Write(v)
}

func (w *noGrow) WriteString(s string) (int, error) {
	return w.w.WriteString(s)
}

type noWriteString struct {
	w *bytes.Buffer
}

func (w *noWriteString) Grow(n int) {
	w.w.Grow(n)
}

func (w *noWriteString) Write(v []byte) (int, error) {
	return w.w.Write(v)
}

type justWriter struct {
	w *bytes.Buffer
}

func (w *justWriter) Write(v []byte) (int, error) {
	return w.w.Write(v)
}

func cvtElems(v []string) (out [][]byte) {
	for _, s := range v {
		out = append(out, []byte(s))
	}
	return out
}

func TestString(t *testing.T) {
	buf := &bytes.Buffer{}

	// test with every combination of optional methods
	writers := []struct {
		Name string
		w    io.Writer
	}{
		{"buffer", buf},
		{"noGrow", &noGrow{w: buf}},
		{"noWriteString", &noWriteString{w: buf}},
		{"justWriter", &justWriter{w: buf}},
	}

	for i, test := range table {
		expect := strings.Join(test.Elems, test.Sep)
		bSep := []byte(test.Sep)
		bElems := cvtElems(test.Elems)
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {

			for _, writer := range writers {
				t.Run(writer.Name, func(t *testing.T) {
					tests := []struct {
						Name string
						Test func(io.Writer) (string, int64, error)
					}{
						{"String", func(w io.Writer) (string, int64, error) {
							n, err := String(w, test.Elems, test.Sep)
							return buf.String(), n, err
						}},
						{"Bytes", func(w io.Writer) (string, int64, error) {
							n, err := Bytes(w, bElems, bSep)
							return buf.String(), n, err
						}},
					}

					for _, testFunc := range tests {
						t.Run(testFunc.Name, func(t *testing.T) {
							s, n, err := testFunc.Test(writer.w)
							if err != nil {
								t.Fatalf("unexpected error: %v", err)
							}
							if int(n) != len(expect) {
								t.Fatalf("expected %d bytes, got %d", n, len(expect))
							}
							if s != expect {
								t.Fatalf("expected %q, got %q", expect, s)
							}
							t.Log(s)
							buf.Reset()
						})
					}
				})
			}
		})
	}
}
