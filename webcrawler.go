package webcrawler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/antchfx/htmlquery"
)

var visited = make(map[string]bool)

func RunCli() {

	if len(os.Args) < 2 {
		fmt.Println("please enter website to crawl for example `go run cmd/main.go https://monzo.com`")
		os.Exit(1)
	}

	CrawlPage(os.Args[1])
}

func CrawlPage(website string) {

	var sm sync.Map
	sm.Store(website, true)
	// visited[website] = true

	links, err := ProcessWebPage(website)
	if err != nil {
		fmt.Println(err)
	}

	for _, link := range links {
		if sm.Load(website); !visited[link] {
			fmt.Printf("crawling %s \n", link)
			go CrawlPage(link)
			continue
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
