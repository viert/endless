package endless

import (
	"testing"
)

var (
	first  []byte = []byte{1, 2, 3, 4, 5}
	second []byte = []byte{6, 7, 8, 9, 10, 11, 12}

	result []byte = []byte{11, 12, 3, 4, 5, 6, 7, 8, 9, 10}

	write1 []byte = []byte{1, 2, 3, 4, 5, 6, 7}
	write2 []byte = []byte{8, 9, 10, 11, 12, 13, 14}
	write3 []byte = []byte{15, 16, 17, 18, 19, 20, 21}

	read1 []byte = []byte{1, 2, 3}
	read2 []byte = []byte{4, 5, 6}
	read3 []byte = []byte{7, 8, 9}
	read4 []byte = []byte{10, 11, 12}
	read5 []byte = []byte{13, 14}
)

func assertItemsEqual(t *testing.T, a []byte, b []byte) {
	if len(a) != len(b) {
		t.Errorf("arrays are of different length. len(a) = %d, len(b) = %d", len(a), len(b))
	}

	for i, item := range a {
		if b[i] != item {
			t.Errorf("wrong item a[%d] = %d, b[%d] = %d", i, item, i, b[i])
		}
	}
}

func TestEndless(t *testing.T) {
	e := NewEndless(10)
	e.Write(first)
	for i, item := range first {
		if e.data[i] != item {
			t.Errorf("unexpected value at %d. expected %d, but got %d", i, item, e.data[i])
		}
	}

	if e.start != 0 {
		t.Errorf("start has changed, expected to be zero")
	}

	if e.pos() != 5 {
		t.Errorf("e.pos() is expected to be 5 but is %d", e.pos())
	}

	if e.Filled() {
		t.Error("Filled() should return false at this point")
	}

	e.Write(second)
	for i, item := range result {
		if e.data[i] != item {
			t.Errorf("unexpected value at %d. expected %d, but got %d", i, item, e.data[i])
		}
	}

	if e.pos() != 2 {
		t.Errorf("e.pos() is expected to be 2 but is %d", e.pos())
	}

	if e.start != 2 {
		t.Errorf("start is expected to be 2 but is %d", e.start)
	}

	if e.writeCursor != uint64(len(first)+len(second)) {
		t.Errorf("writeCursor must be at position %d but is at %d instead", len(first)+len(second), e.writeCursor)
	}

	e.Write(write1)
	if !e.Filled() {
		t.Error("Filled() should return true at this point")
	}
}

func TestReader(t *testing.T) {
	var n int
	var err error

	e := NewEndless(10)
	r := e.NewReader(0)

	buf := make([]byte, 3)

	e.Write(write1)
	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 3 {
		t.Errorf("expected to read 3 items but read %d", n)
	}
	assertItemsEqual(t, buf, read1)

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 3 {
		t.Errorf("expected to read 3 items but read %d", n)
	}
	assertItemsEqual(t, buf, read2)

	e.Write(write2)

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 3 {
		t.Errorf("expected to read 3 items but read %d", n)
	}
	assertItemsEqual(t, buf, read3)

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 3 {
		t.Errorf("expected to read 3 items but read %d", n)
	}
	assertItemsEqual(t, buf, read4)

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 2 {
		t.Errorf("expected to read 2 items but read %d", n)
	}
	assertItemsEqual(t, buf[:n], read5)

	n, err = r.Read(buf)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if n != 0 {
		t.Errorf("expected to read 0 items but read %d", n)
	}

	e.Write(write3)
	e.Write(write1)
	e.Write(write2)

	n, err = r.Read(buf)
	if err == nil {
		t.Errorf("should have given an error that the reader has remained behind")
	}

}
