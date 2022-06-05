# webcrawler

to run

`go run cmd/main.go https://monzo.com`

to test

`go test`

#TODO

- I would like to write a handful of small HTML pages that link to each other, and test that crawling from the root page traverses the other pages.  At the moment the tests rely on traversing a single internal html page.  That is not really a demanding test, also it really makes no effort to test our concurrency model.

- I dont like how concurrency is configured here, at the moment concurrency is hardcoded to 5

`	for i := 0; i < 5; i++ {
		go func() {
`

  A more elegant solution would be to have something like

`	for _, link := range links {
		site, ok := sm.Load(link)
		if !ok {
			fmt.Printf("crawling %s \n", site)
			go CrawlPage(site.(string))
		}
`

  Please note this is not a working example we would need to think about how we deal with the global shared data (in this case I think a   sync.map might be a good option).  This approach would allow the number of workers to dynamically grow as the number of links grow.

- I would like to check the HTTP response status and do something reasonable is a non 200 is returned, also I would like to optimise the crawl further by checking for media type (no point trying to parse MP3s for links, for example).

