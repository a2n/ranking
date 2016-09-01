package ranking

import (
	"net/http"
	"encoding/json"
	"net/url"
	"io/ioutil"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
)

type Bing struct {
	key string
}

type BingLocale int
const (
	es_AR BingLocale = iota
	en_AU
	de_AT
	nl_BE
	fr_BE
	pt_BR
	en_CA
	fr_CA
	es_CL
	da_DK
	fi_FI
	fr_FR
	de_DE
	zh_HK
	en_IN
	en_ID
	en_IE
	it_IT
	ja_JP
	ko_KR
	en_MY
	es_MX
	nl_NL
	en_NZ
	no_NO
	zh_CN
	pl_PL
	pt_PT
	en_PH
	ru_RU
	ar_SA
	en_ZA
	es_ES
	sv_SE
	fr_CH
	de_CH
	zh_TW
	tr_TR
	en_GB
	en_US
	es_US
)

func NewBing(key string) (*Bing) {
	if len(key) == 0 {
		panic("No subscription key.")
	}

	return &Bing{
		key: key,
	}
}

func (this *Bing) Get(locale BingLocale, host, kw string, max int) (int, error) {
	req, e := this.request(locale, kw, max)
	if e != nil {
		return notFound, errors.Wrap(e, "Get failed")
	}

	resp, e := http.DefaultClient.Do(req)
	if e != nil {
		return notFound, errors.Wrap(e, "Get failed")
	}

	b, e := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if e != nil {
		return notFound, errors.Wrap(e, "Get failed")
	}

	return this.parse(host, b)
}

func (this *Bing) request(locale BingLocale, kw string, max int) (*http.Request, error) {
	loc, e := this.localeText(locale)
	if e != nil {
		return nil, errors.Wrap(e, "create request failed")
	}

	if len(kw) == 0 {
		return nil, errors.New("create request failed, empty keyword.")
	}

	urlstr := fmt.Sprintf("https://api.cognitive.microsoft.com/bing/v5.0/search?count=%d&mkt=%s&q=%s", max * 10, loc, url.QueryEscape(kw))

	req, e := http.NewRequest("GET", urlstr, nil)
	if e != nil {
		return nil, errors.Wrap(e, "create request failed")
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", this.key)

	return req, nil
}

func (this *Bing) localeText(l BingLocale) (string, error) {
	hash := map[BingLocale]string {
		es_AR: "es-AR",
		en_AU: "en-AU",
		de_AT: "de-AT",
		nl_BE: "nl-BE",
		fr_BE: "fr-BE",
		pt_BR: "pt-BR",
		en_CA: "en-CA",
		fr_CA: "fr-CA",
		es_CL: "es-CL",
		da_DK: "da-DK",
		fi_FI: "fi-FI",
		fr_FR: "fr-FR",
		de_DE: "de-DE",
		zh_HK: "zh-HK",
		en_IN: "en-IN",
		en_ID: "en-ID",
		en_IE: "en-IE",
		it_IT: "it-IT",
		ja_JP: "ja-JP",
		ko_KR: "ko-KR",
		en_MY: "en-MY",
		es_MX: "es-MX",
		nl_NL: "nl-NL",
		en_NZ: "en-NZ",
		no_NO: "no-NO",
		zh_CN: "zh-CN",
		pl_PL: "pl-PL",
		pt_PT: "pt-PT",
		en_PH: "en-PH",
		ru_RU: "ru-RU",
		ar_SA: "ar-SA",
		en_ZA: "en-ZA",
		es_ES: "es-ES",
		sv_SE: "sv-SE",
		fr_CH: "fr-CH",
		de_CH: "de-CH",
		zh_TW: "zh-TW",
		tr_TR: "tr-TR",
		en_GB: "en-GB",
		en_US: "en-US",
		es_US: "es-US",
	}

	s, ok := hash[l]
	if !ok {
		return "", errors.New("Invalid locale")
	}
	return s, nil
}

type bingGetResult struct {
	Pages struct {
		Value []struct {
			Name string `json:"name"`
			URL string `json:"displayUrl"`
		} `json:"value"`
	} `json:"webPages"`
}

func (this *Bing) parse(host string, b []byte) (int, error) {
	if len(b) == 0 {
		return notFound, errors.New("parse failed, empty content.")
	}

	if len(host) == 0 {
		return notFound, errors.New("parse failed, empty host.")
	}

	var result bingGetResult
	if e := json.Unmarshal(b, &result); e != nil {
		return notFound, errors.Wrap(e, "Get failed")
	}

	for k, item := range result.Pages.Value {
		log.Printf("%d, %s", k, item.URL)

		if strings.Index(item.URL, host) != -1 {
			return k, nil
		}
	}

	return notFound, nil
}
