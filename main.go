package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github/Sitemap/parse"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build a sitemap for")
	maxDepth := flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXML := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()
	return filter(base, hrefs(resp.Body, base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := parse.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func filter(b string, st []string) []string {
	var ret []string
	for _, text := range st {
		if strings.HasPrefix(text, b) {
			ret = append(ret, text)
		}

	}
	return ret
}

func bfs(urlStr string, depth int) []string {
	seen := make(map[string]struct{})

	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i <= depth; i++ {
		q, nq := nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}
		for urls := range q {
			if _, ok := seen[urls]; ok {
				continue
			}
			seen[urls] = struct{}{}
			for _, link := range get(urls) {
				nq[link] = struct{}{}
			}

		}

	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret

}
