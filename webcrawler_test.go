package webcrawler_test

import (
	"fmt"
	"io/ioutil"
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
