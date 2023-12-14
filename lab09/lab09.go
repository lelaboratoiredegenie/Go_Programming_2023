package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
)

func main() {

	maxFlag := flag.Int("max", 10, "Max number of comments to show")
	flag.Parse()
	c := colly.NewCollector()
	idx := 1

	c.OnHTML(".push", func(e *colly.HTMLElement) {
		if *maxFlag > 0 {
			userID := e.ChildText(".push-userid")
			content := e.ChildText(".push-content")
			time := strings.TrimSpace(e.ChildText(".push-ipdatetime"))

			fmt.Printf("%d. 名字：%s，留言%s，時間： %s\n", idx, userID, content, time)
			*maxFlag--
			idx++
		}
	})

	url := "https://www.ptt.cc/bbs/joke/M.1481217639.A.4DF.html"
	c.Visit(url)
}
