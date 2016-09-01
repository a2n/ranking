package ranking

import "testing"

func TestGoogle(t *testing.T) {
	n, e := NewGoogle().Get("gov.tw", "Taiwan", 5)
	if e != nil {
		panic(e)
	}

	t.Logf("Index: %d", n)
}
