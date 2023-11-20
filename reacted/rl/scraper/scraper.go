package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

func FetchTitle(url string) (string, error) {
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	title := doc.Find("title").Text()

	return title, nil
}