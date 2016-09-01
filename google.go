package ranking

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"strconv"
	"regexp"
	"bytes"
	"html"
	"log"
	"io"

	"github.com/pkg/errors"
)

type Google struct {
	page *regexp.Regexp
	pageReplace *regexp.Regexp

	host *regexp.Regexp
	markup *regexp.Regexp

	max int
	nowPage int
	items int
}

func NewGoogle() (*Google) {
	return &Google {
		page: regexp.MustCompile(`/search[^"]+start=\d+[^"]+`),
		pageReplace: regexp.MustCompile(`start=\d+`),
		host: regexp.MustCompile(`<cite[^>]*>([^<]+)</cite>`),
		markup: regexp.MustCompile(`<[^>]*>`),
	}
}

func (this *Google) request(urlstr string, cookies []*http.Cookie) (*http.Request, error) {
	r, e := http.NewRequest("GET", urlstr, nil)
	if e != nil {
		return nil, errors.Wrap(e, "get request failed")
	}

	// User agent
	//r.Header = UserAgent(r.Header)

	// Cookies
	for _, c := range cookies {
		r.AddCookie(c)
	}

	return r, nil
}

func (this *Google) Get(host, kw string, max int) (int, error) {
	if len(host) == 0 {
		return 0, errors.New("Get failed, empty host.")
	}

	if len(kw) == 0 {
		return 0, errors.New("Get failed, empty keywords.")
	}

	this.max = max
	const base = "https://www.google.com/search?q="
	resp, e := http.DefaultClient.Get(base + url.QueryEscape(kw))
	if e != nil {
		resp.Body.Close()
		return 0, errors.Wrap(e, "get failed")
	}

	this.nowPage = 0
	n, e := this.get(host, resp)
	if e != nil {
		log.Fatalf("%+v", e)
	}

	return n, nil
}

func (this *Google) get(host string, resp *http.Response) (int, error) {
	//log.Printf("resp: %s", resp.Request.URL.String())

	b, e := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if e != nil {
		return notFound, errors.Wrap(e, "get failed")
	}

	// Hosts
	if n := this.findHost(host, b); n != notFound {
		return this.items, nil
	}

	if this.nowPage + 1 >= this.max {
		return notFound, io.EOF
	}
	this.nowPage++

	// Pages
	urlstr := "https://www.google.com.tw" + this.findPage(b)
	req, e := this.request(urlstr, resp.Cookies())
	if e != nil {
		return notFound, errors.Wrap(e, "get failed")
	}

	resp, e = http.DefaultClient.Do(req)
	if e != nil {
		return notFound, errors.Wrap(e, "get failed")
	}
	return this.get(host, resp)
}

func (this *Google) findHost(host string, b []byte) (int) {
	hosts := this.host.FindAllSubmatch(b, -1)
	if len(hosts) == 0 {
		log.Print("host not found.")
		return notFound
	}

	for _, v := range hosts {
		this.items++

		v1 := this.markup.ReplaceAll(v[0], []byte{})
		log.Printf("%d  %s", this.items, v1)
		n := strings.Index(string(v1), host)
		if n != -1 {
			return this.items
		}
	}
	return notFound
}

func (this *Google) findPage(b []byte) string {
	if len(b) == 0 {
		return ""
	}

	ms := this.page.FindAll(b, -1)
	if len(ms) == 0 {
		return ""
	}

	var target string
	next := []byte(strconv.Itoa((this.nowPage + 1) * 10))
	for _, m := range ms {
		if bytes.Index(m, next) != -1 {
			target = string(m)
		}
	}

	return html.UnescapeString(target)
}
