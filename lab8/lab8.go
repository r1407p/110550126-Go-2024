package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"log"
)

func main() {
	maxComments := flag.Int("max", 10, "Max number of comments to show")
	flag.Parse()

	c := colly.NewCollector()
	commentCount := 0

	c.OnHTML(".push", func(e *colly.HTMLElement) {
		if commentCount >= *maxComments {
			return
		}
	
		userID := e.ChildText(".push-userid")
		content := e.ChildText(".push-content")
		time := e.ChildText(".push-ipdatetime")
		if userID != "" && content != "" {
			fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", commentCount+1, userID, content, time)
			commentCount++
		}
	})

	url := "https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html"
	if err := c.Visit(url); err != nil {
		log.Fatalf("error visiting url: %v", err)
	}
	
}
