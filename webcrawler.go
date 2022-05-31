package webcrawler

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
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

	linkChan := make(chan []string)
	notVisitedChan := make(chan string)

	go func() {
		linkChan <- websites
	}()

	go func() {
		for link := range notVisitedChan {
			links, err := ProcessWebPage(link)
			if err != nil {
				fmt.Println(err)
			}
			go func() {
				linkChan <- links
			}()
			fmt.Printf("crawling %s \n", link)
		}
	}()

	visited := make(map[string]bool)
	for links := range linkChan {
		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				notVisitedChan <- link
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

	links, err := FindUrls(content)
	if err != nil {
		return nil, err
	}

	urls, err := Canonicalise(links, url)
	if err != nil {
		return nil, err
	}

	return urls, err
}

func Crawl(website string) (*http.Response, error) {

	resp, err := http.Get(website)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func FindUrls(resp *http.Response) ([]string, error) {

	var links []string

	node, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlNodes, _ := htmlquery.QueryAll(node, "//a/@href")
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

func Canonicalise(links []string, url *url.URL) ([]string, error) {

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
