package scraper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"log"
)

const(
	pageUrl = "https://mercatox.com/exchange/VERI/ETH"
)

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("_examples/challenge.html")
		if err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Server", "cloudflare")
		w.WriteHeader(503)
		w.Write(b)
	}))
	defer ts.Close()

	scraper, err := NewTransport(http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	c := http.Client{
		Transport: scraper,
	}

	u, err := url.Parse(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	var cookie string
	for {
		if cookie != "" {
			break
		}

		log.Printf("Getting the cookie")

		_, err := c.Get(u.String())
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range scraper.GetCookies().Cookies(u) {
			if c.Name == "cf_clearance" {
				cookie = c.Value
				break
			}
		}
	}

	fmt.Printf("cf_clearance=%s", cookie)

	res, err := c.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}
