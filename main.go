package main

import (
	"flag"
	"github/Sitemap/parse"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	flagurl := flag.String("url", "http://gophercises.com/", " the url of the site for which sitemap is to be made")
	flag.Parse()
	resp, err := http.Get(*flagurl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	req := resp.Request.URL
	baseurl := &url.URL{
		Scheme: req.Scheme,
		Host:   req.Host,
	}
	base := baseurl.String()
	links, err := parse.Parse(resp.Body)
	var note []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			note = append(note, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			note = append(note, l.Href)

		}
	}

}
