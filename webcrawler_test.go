package webcrawler_test

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"testing"
	"webcrawler"

	"github.com/google/go-cmp/cmp"
)

func TestCrawl(t *testing.T) {
	got, err := webcrawler.Crawl("https://www.example.com/")
	if err != nil {
		fmt.Println(err)
	}

	want, err := ioutil.ReadFile("testdata/webcrawler.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestCanonicalise(t *testing.T) {

	url, err := url.Parse("https://monzo.com")
	if err != nil {
		t.Fatal(err)
	}

	links := []string{"/i/business", "/i/current-account/", "https://monzo.com/faq", "https://app.adjust.com/ydi27sn_9mq4ox7?engagement_type=fallback_click&fallback=https://monzo.com/download&redirect_macos=https://monzo.com/download"}

	got, err := webcrawler.Canonicalise(links, url)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"https://monzo.com/i/business", "https://monzo.com/i/current-account/", "https://monzo.com/faq"}

	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
