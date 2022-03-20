package webcrawler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func RunCli() {

	links, _ := ProcessWebPage(os.Args[1])

	for _, link := range links {
		fmt.Println(link)
		ProcessWebPage(link)
	}

}

func ProcessWebPage(website string) ([]string, error) {

	url, err := url.Parse(website)
	if err != nil {
		log.Fatal(err)
	}
	content, err := Crawl(website)
	if err != nil {
		log.Fatal(err)
	}

	links, err := findUrls(url, content)
	if err != nil {
		log.Fatal(err)
	}
	urls, err := uniquePaths(links, url)
	if err != nil {
		fmt.Println(err)
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

	return content, nil
}

func findUrls(urlToGet *url.URL, content []byte) ([]string, error) {

	var (
		err       error
		links     []string = make([]string, 0)
		findLinks          = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	)

	// Retrieve all anchor tag URLs from string
	matches := findLinks.FindAllStringSubmatch(string(content), -1)

	for _, val := range matches {
		var linkUrl *url.URL

		// Parse the anchr tag URL
		if linkUrl, err = url.Parse(val[1]); err != nil {
			return links, err
		}

		// If the URL is absolute, add it to the slice
		// If the URL is relative, build an absolute URL
		if linkUrl.IsAbs() {
			links = append(links, linkUrl.String())
		} else {
			links = append(links, urlToGet.Scheme+"://"+urlToGet.Host+linkUrl.String())
		}
	}

	return links, err
}

func uniquePaths(links []string, url *url.URL) ([]string, error) {

	var paths []string

	for _, v := range links {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		if u.Host == url.Host {
			paths = append(paths, u.Path)
		}
	}

	inResult := make(map[string]bool)
	var uniquePaths []string
	for _, str := range paths {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			if str != "/" {
				if str != "" {
					uniquePaths = append(uniquePaths, "https://"+url.Host+str)
				}
			}
		}
	}
	return uniquePaths, nil
}
