package webcrawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
)

func RunCli() {

	if len(os.Args) < 2 {
		fmt.Println("please enter website to crawl for example `go run cmd/main.go https://monzo.com`")
		os.Exit(1)
	}

	CrawlPage(os.Args[1])
}

func CrawlPage(website string) {

	websites := []string{website}

	worklist := make(chan []string)  // lists of URLs, may have duplicates
	unseenLinks := make(chan string) // de-duplicated URLs

	go func() { worklist <- websites }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks, _ := ProcessWebPage(link)
				go func() { worklist <- foundLinks }()
				fmt.Println(link)
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				unseenLinks <- link
			}
		}
	}
}

func ProcessWebPage(website string) ([]string, error) {

	url, err := url.Parse(website)
	if err != nil {
		return nil, err
	}

	content, err := Crawl(website)
	if err != nil {
		return nil, err
	}

	links, err := findUrls(url, content)
	if err != nil {
		return nil, err
	}

	urls, err := canonicalise(links, url)
	if err != nil {
		return nil, err
	}

	return urls, err
}

func Crawl(website string) ([]byte, error) {

	resp, err := http.Get(website)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return content, err
}

func findUrls(urlToGet *url.URL, content []byte) ([]string, error) {

	var links []string

	doc, err := htmlquery.LoadURL(urlToGet.String())
	if err != nil {
		return nil, err
	}
	htmlNodes, _ := htmlquery.QueryAll(doc, "//a/@href")
	for _, n := range htmlNodes {
		href := htmlquery.SelectAttr(n, "href")
		url, err := url.Parse(href)
		if err != nil {
			return nil, err
		}
		links = append(links, url.String())
	}
	return links, err
}

func canonicalise(links []string, url *url.URL) ([]string, error) {

	var matchingLinks []string

	for _, l := range links {
		link, err := url.Parse(l)
		if err != nil {
			return nil, err
		}
		if link.IsAbs() && link.Host == url.Host {
			matchingLinks = append(matchingLinks, link.String())
		}
		if strings.HasPrefix(link.String(), "/") {
			matchingLinks = append(matchingLinks, link.Scheme+"://"+link.Host+link.String())
		}
	}

	return matchingLinks, nil
}
