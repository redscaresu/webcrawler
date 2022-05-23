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

var visited = make(map[string]bool)

func RunCli() {

	CrawlPage(os.Args[1])

}

func CrawlPage(website string) {

	visited[website] = true

	links, err := ProcessWebPage(website)
	if err != nil {
		fmt.Println(err)
	}

	for _, link := range links {
		if !visited[link] {
			// time.Sleep(1 * time.Second)
			fmt.Printf("crawling %s \n", link)
			CrawlPage(link)
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
		u, err := url.Parse(href)
		if err != nil {
			return nil, err
		}
		if u.IsAbs() && u.Host == urlToGet.Host {
			links = append(links, u.String())
		}
		if strings.HasPrefix(u.String(), "/") {
			links = append(links, urlToGet.Scheme+"://"+urlToGet.Host+u.String())
		}
	}
	return links, err
}

func canonicalise(links []string, url *url.URL) ([]string, error) {

	var paths []string
	// var doNotFollow []string

	for _, v := range links {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		if u.Host == url.Host {
			paths = append(paths, u.Path)
			// } else {
			// 	doNotFollow = append(doNotFollow, u.String())
			// }
		}
	}

	var uniquePaths []string
	for _, str := range paths {
		if str != "/" || str != "" {
			uniquePaths = append(uniquePaths, "https://"+url.Host+str)
		}
	}

	return uniquePaths, nil
}
