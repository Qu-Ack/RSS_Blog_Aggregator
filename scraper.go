package main

import (
	"context"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Atom    string   `xml:"atom,attr"`
	Channel struct {
		Text  string `xml:",chardata"`
		Title string `xml:"title"`
		Link  struct {
			Text string `xml:",chardata"`
			Href string `xml:"href,attr"`
			Rel  string `xml:"rel,attr"`
			Type string `xml:"type,attr"`
		} `xml:"link"`
		Description   string `xml:"description"`
		Generator     string `xml:"generator"`
		Language      string `xml:"language"`
		LastBuildDate string `xml:"lastBuildDate"`
		Item          []struct {
			Text        string `xml:",chardata"`
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			PubDate     string `xml:"pubDate"`
			Guid        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

func fetchXMLfromFEED(url string) (Rss, error) {
	res, err := http.Get(url)
	if err != nil {
		return Rss{}, err
	}

	body, err := io.ReadAll(res.Body)

	xml_data := Rss{}
	err = xml.Unmarshal(body, &xml_data)

	if err != nil {
		return Rss{}, err
	}

	return xml_data, nil

}

func (cfg apiConfig) scraper() {
	ticker := time.NewTicker(60 * time.Second)

	defer ticker.Stop()

	for ; ; <-ticker.C {
		feeds, err := cfg.getNextFeedsToFetch(context.Background())
		if err != nil {
			return
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				xml_data, err := fetchXMLfromFEED(feed.Url)

				if err != nil {
					log.Println("Can't get xml")
					return
				}

				cfg.markFeedFetched(context.Background(), feed.ID)

				for _, item := range xml_data.Channel.Item {
					log.Println(item.Title)
				}

				log.Printf("Feed %s collected, %v posts found", feed.Name, len(xml_data.Channel.Item))

			}(&wg)

		}
		wg.Wait()

	}
}
