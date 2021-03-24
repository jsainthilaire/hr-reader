package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Data struct {
	Items []Item
}

type Item struct {
	Id            string
	Title         string
	ContentUrl    string `json:"content_html"`
	Url           string
	ExternalUrl   string `json:"external_url"`
	DatePublished string `json:"date_published"`
	Author        string
	Description   string
	ImageURL      string
}

const key = "key from link preview"
const hnFeedUrl = "https://hnrss.org/frontpage.jsonfeed"
const hnRawDataFile = "hn.json"
const hnFullDataFile = "hnfull.json"

func fetchData() error {
	resp, err := http.Get(hnFeedUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	f, err := os.Create(hnRawDataFile)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := fetchLinkPreviewInfo()
	if err != nil {
		log.Fatalln(err)
	}
}

func fetchLinkPreviewInfo() error {
	f, err := ioutil.ReadFile(hnRawDataFile)
	if err != nil {
		return err
	}

	data := Data{}
	err = json.Unmarshal(f, &data)
	if err != nil {
		return err
	}

	var items []Item
	for _, item := range data.Items {
		// TODO: handle error
		lp, _ := linkPreview(item.Url)

		item.Description = lp.Description
		item.ImageURL = lp.Image

		items = append(items, item)
	}

	data.Items = items
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(hnFullDataFile, b, 0644)
}

type LinkPreviewItem struct {
	Image       string
	Description string
}

func linkPreview(url string) (LinkPreviewItem, error) {
	fullUrl := fmt.Sprintf("%s?key=%s&q=%s", "http://api.linkpreview.net/", key, url)
	resp, err := http.Get(fullUrl)
	if err != nil {
		return LinkPreviewItem{}, nil
	}

	defer resp.Body.Close()

	// TODO: handle error
	b, _ := ioutil.ReadAll(resp.Body)
	data := LinkPreviewItem{}
	json.Unmarshal(b, &data)

	return data, nil
}
