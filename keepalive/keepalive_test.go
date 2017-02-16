package keepalive

import "testing"

func TestDoIt(t *testing.T) {
	size := int64(25000)
	s := make([]int64, size)

	var want int64
	for i := int64(0); i < size; i++ {
		s[i] = i
		want += i
	}

	got := DoIt(s)

	if want != got {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}

func TestDoItKeepAlive(t *testing.T) {
	size := int64(25000)
	s := make([]int64, size)

	var want int64
	for i := int64(0); i < size; i++ {
		s[i] = i
		want += i
	}

	got := DoItKeepAlive(s)

	if want != got {
		t.Errorf("Wanted %d, got %d", want, got)
	}
}
