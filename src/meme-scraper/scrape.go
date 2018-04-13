package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anaskhan96/soup"
	//"golang.org/x/net/html"
	//	"golang.org/x/net/html/atom"
)

type meme struct {
	name string
	src  string
}

func mainPage(page int) ([]meme, error) {
	pageURL := fmt.Sprintf("%smemes/page/%d", *url, page)
	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Non-200 response code received: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	doc := soup.HTMLParse(newStr)

	rows := doc.FindAll("tr")
	toScrape = []meme{}
	for _, r := range rows {
		elements := r.FindAll("td")
		for _, element := range elements {
			img := element.Find("img")
			a := element.Find("a")
			if img.Error != nil || a.Error != nil {
				continue
			}
			meme := meme{}
			for _, attr := range img.Pointer.Attr {
				if attr.Key == "title" {
					meme.name = attr.Val
				}
			}
			if meme.name == "" {
				continue
			}

			for _, attr := range a.Pointer.Attr {
				if attr.Key == "href" {
					meme.src = attr.Val
				}
			}
			fmt.Println(meme)
			toScrape = append(toScrape, meme)
		}
	}

	return toScrape, nil
}

func scrapeMemePages()

func scrape() {
	page := 1
	for {
		time.Sleep(time.Millisecond * time.Duration(*rate))
		toScrape, err := mainPage(page)
		if err != nil {
			// KISS and retry
			log.Printf("Could not get page %d: %v", page, err)
			break
		}

		scrapeMemePages()
		page++
	}
}
