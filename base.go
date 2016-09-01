package ranking

import (
	"log"
	"net/http"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}

func UserAgent(h http.Header) (http.Header) {
	// Clone
	h1 := make(http.Header)
	for k, v := range h {
		h1.Add(k, v[0])
	}

	h1.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36") 

	return h1
}

const (
	notFound = -1
)
