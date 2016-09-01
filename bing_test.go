package ranking

import (
	"testing"
)

func TestBing(t *testing.T) {
	s := NewBing("")
	n, e := s.Get(zh_TW, "gov.tw", "Taiwan", 5)
	if e != nil {
		t.Fatal(e)
	}

	t.Logf("Index: %v\n", n)
}
