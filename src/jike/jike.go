package jike

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func main() {
	var mainUrl string
	mainUrl = "https://time.geekbang.org/library?category=0"
	//mainUrl = "https://time.geekbang.org/library?category=1"
	//启动一个浏览器
	content := StartChromeAndGetContent(mainUrl)
	log.Println(content)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		log.Fatal(" goquery.NewDocumentFromReader error:", err)
	}

	total := 0
	doc.Find("ul.column-list li a").Each(func(i int, s *goquery.Selection) {
		// url
		if v, exist := s.Attr("href"); exist {
			log.Println(v)
		}
		total += 1
	})
	log.Println("total", total)
}
