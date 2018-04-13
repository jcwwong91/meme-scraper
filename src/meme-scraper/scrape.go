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

func mainPage(page int) error {
	pageURL := fmt.Sprintf("%smemes/page/%d", *url, page)
	resp, err := http.Get(pageURL)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Non-200 response code received: %d", resp.StatusCode)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	doc := soup.HTMLParse(newStr)

	rows := doc.FindAll("tr")
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
		}
	}

	return nil
}

func scrape() {
	page := 1
	for {
		time.Sleep(time.Millisecond * time.Duration(*rate))
		err := mainPage(page)
		if err != nil {
			// KISS and retry
			log.Printf("Could not get page %d: %v", page, err)
			break
		}
		page++
	}
}
