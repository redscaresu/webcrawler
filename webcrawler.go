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

	links, _, err := ProcessWebPage(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	for _, link := range links {
		ManageCrawlers(link)
	}
}

func ManageCrawlers(link string) {
	visitedLinks := make(map[string]bool)

	links, _, err := ProcessWebPage(link)
	if err != nil {
		fmt.Println(err)
	}
	visitedLinks[os.Args[1]] = true

	for _, link := range links {
		visitedLinks[link] = false
	}

	for k, _ := range visitedLinks {
		if !visitedLinks[k] {
			visitedLinks[link] = true
			links, _, _ := ProcessWebPage(k)
			for _, link := range links {
				visitedLinks[link] = true
				fmt.Printf("crawled page: %v link: %s\n", k, link)
			}
		}
	}
}

func ProcessWebPage(website string) ([]string, []string, error) {

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
	urls, doNotFollow, err := uniquePaths(links, url)
	if err != nil {
		fmt.Println(err)
	}

	return urls, doNotFollow, err
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

func uniquePaths(links []string, url *url.URL) ([]string, []string, error) {

	var paths []string
	var doNotFollow []string

	for _, v := range links {
		u, err := url.Parse(v)
		if err != nil {
			return nil, nil, err
		}
		if u.Host == url.Host {
			paths = append(paths, u.Path)
		} else {
			doNotFollow = append(doNotFollow, u.String())
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

	for _, v := range doNotFollow {
		fmt.Printf("do not follow %s\n", v)
	}

	return uniquePaths, doNotFollow, nil
}
