package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/deckarep/golang-set"
)

const PPTUrl = "http://atwoodknives.blogspot.nl/"

var latestSet = mapset.NewSet()
var (
	pushedAppSecret string
	pushedAppKey    string
)

func main() {

	pushedAppSecret = os.Getenv("PUSHEDSECRET")
	pushedAppKey = os.Getenv("PUSHEDAPPKEY")

	latestSet = getCurrentSetOfPostLinks()
	sendPush(fmt.Sprintf("%s", latestSet.ToSlice()[0]))

	//check every 5 seconds..
	doEvery(5000*time.Millisecond, calculatePostSetDifference)
	//fmt.Printf("New set of post urls of size %d is:", len(postLinks.ToSlice()))
	//for link := range postLinks.Iter() {
	//fmt.Printf("%s\n", link)
	//}

	//we need //*[@id="Blog1"]/div[1]
	//div class blog-posts

	//calculate a hash/checksum
	//compare to last known one in the database
	//if different send push notification
}

func sendPush(payload string) (resp *http.Response, err error) {
	return http.PostForm("https://api.pushed.co/1/push",
		url.Values{"app_key": {pushedAppKey},
			"app_secret":    {pushedAppSecret},
			"target_type":   {"app"},
			"content_type":  {"url"},
			"content_extra": {payload},
			"content":       {payload},
		})
}

func calculatePostSetDifference(t time.Time) {
	newSet := getCurrentSetOfPostLinks()
	if newSet.IsSubset(latestSet) {
		log.Printf("No new items found, found %d items\n", len(newSet.ToSlice()))
	} else {
		newItems := newSet.Difference(latestSet)
		latestSet = newSet
		//notify....
		log.Print("BUY BUY BUY")
		sendPush(fmt.Sprintf("BUY NOW: %s", newItems.ToSlice()[0]))
	}
}

func getCurrentSetOfPostLinks() mapset.Set {
	doc, err := goquery.NewDocument("http://atwoodknives.blogspot.nl")
	if err != nil {
		log.Fatal(err)
	}

	//#Blog1 > div.blog-posts.hfeed > div:nth-child(1) > div > div > div > h3 > a
	postLinks := mapset.NewSet()
	doc.Find(".blog-posts").Each(func(i int, s *goquery.Selection) {
		links := s.Find("a")
		links.Each(func(k int, link *goquery.Selection) {
			linkText, _ := link.Attr("href")
			if linkText != "" && strings.Contains(linkText, "atwoodknives.blogspot") {
				postLinks.Add(linkText)
			}
		})

		postBodies := s.Find(".post-body")
		postBodies.Each(func(j int, body *goquery.Selection) {
			//fmt.Printf("Blog %d: body: %s\n", j, body.Text())
		})
	})
	return postLinks
}

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}
