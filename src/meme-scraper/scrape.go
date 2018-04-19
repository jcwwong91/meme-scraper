package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"golang.org/x/net/html"
	//	"golang.org/x/net/html/atom"
)

type meme struct {
	Name        string
	Src         string
	Views       int
	Videos      int
	Images      int
	Comments    int
	Created     time.Time
	LastUpdated time.Time
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
	memes := []meme{}
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
					meme.Name = attr.Val
				}
			}
			if meme.Name == "" {
				continue
			}

			for _, attr := range a.Pointer.Attr {
				if attr.Key == "href" {
					meme.Src = attr.Val
				}
			}
			memes = append(memes, meme)
		}
	}

	return memes, nil
}

func parseBasicInfo(attrs []html.Attribute) (int, error) {
	str := strings.Split(attrs[1].Val, " ")
	if len(str) != 2 {
		return 0, fmt.Errorf("Invalid format for %s", attrs[0].Val)
	}
	return strconv.Atoi(strings.Replace(str[0], ",", "", -1))

}

func parseTime(node soup.Root) (time.Time, error) {
	attrs := node.Pointer.Attr
	if len(attrs) != 2 {
		return time.Now(), fmt.Errorf("Invalid time format")
	}
	if attrs[0].Key != "class" || attrs[0].Val != "timeago" || attrs[1].Key != "title" {
		return time.Now(), fmt.Errorf("Invalid time format")
	}
	return time.Parse(time.RFC3339, attrs[1].Val)
}

func scrapeMemePages(memes []meme) error {
	for _, v := range memes {
		if memeExists(v.Name) {
			log.Println(v.Name, "already exists in the database, skipping")
			continue
		}
		time.Sleep(time.Millisecond * time.Duration(*rate))
		pageURL := fmt.Sprintf("%s%s", *url, v.Src)
		resp, err := http.Get(pageURL)
		if err != nil {
			// If we fail, we ignore this meme
			log.Printf("Failed to get page information for %s: %v", v.Name, err)
			continue
		}
		if resp.StatusCode != 200 {
			// If we fail, we ignore this meme
			log.Printf("Bad Status %d for %s", resp.StatusCode, v.Name)
			continue
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()
		doc := soup.HTMLParse(newStr)
		dl := doc.Find("dl")
		if dl.Error != nil {
			log.Printf("Failed to find information for %s", v.Name)
			continue
		}
		basicInfo := dl.FindAll("dd")
		for _, bi := range basicInfo {
			attrs := bi.Pointer.Attr
			if len(attrs) != 2 || attrs[0].Key != "class" || attrs[1].Key != "title" {
				log.Printf("Basic info format incorrect for %s", v.Name)
				continue
			}
			switch attrs[0].Val {
			case "views":
				v.Views, err = parseBasicInfo(attrs)
			case "videos":
				v.Videos, err = parseBasicInfo(attrs)
			case "photos":
				v.Images, err = parseBasicInfo(attrs)
			case "comments":
				v.Comments, err = parseBasicInfo(attrs)
			default:
				log.Printf("Unknown basic info '%s' for %s", attrs[0].Val, v.Name)
				continue
			}
			if err != nil {
				log.Printf("Invalid value for %s (%s) for %s: %v", attrs[0].Val, attrs[1].Val, v.Name, err)
				continue
			}
		}

		timeInfo := doc.FindAll("abbr")
		if len(timeInfo) < 2 {
			log.Printf("Invalid time format detected for %s", v.Name)
			continue
		}
		v.Created, err = parseTime(timeInfo[1])
		if err != nil {
			log.Printf("Error parsing first added for %s: %v", v.Name, err)
			continue
		}
		v.LastUpdated, err = parseTime(timeInfo[0])
		if err != nil {
			log.Printf("Error parsing last updated for %s: %v", v.Name, err)
			continue
		}
		saveChan <- v

	}
	return nil
}

func scrape(start int) {
	page := start
	for {
		time.Sleep(time.Millisecond * (time.Duration(*rate) + time.Duration(rand.Int63n(*random))))
		memes, err := mainPage(page)
		if err != nil {
			// KISS and retry
			log.Printf("Could not get page %d: %v", page, err)
			break
		}

		scrapeMemePages(memes)
		page++
	}
}
