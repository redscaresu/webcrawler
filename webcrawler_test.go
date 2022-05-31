package webcrawler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"webcrawler"

	"github.com/google/go-cmp/cmp"
)

func TestResponse(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("testdata/webcrawler.txt")
		if err != nil {
			t.Fatal(err)
		}
		io.Copy(w, file)
	}))
	defer ts.Close()

	client := ts.Client()
	res, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	got, err := webcrawler.FindUrls(res)
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"https://www.happydoggo.org/domains/example"}
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
