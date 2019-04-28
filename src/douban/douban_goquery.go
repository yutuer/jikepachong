package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strconv"
)

func GetMovie(url string) {
	log.Println("url:", url)

	resp, err := http.Get(url);
	if err != nil {
		log.Fatal(" http.Get error:", err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body);
	if err != nil {
		log.Fatal(" goquery.NewDocumentFromReader error:", err)
	}

	doc.Find("#content h1").Each(func(i int, s *goquery.Selection) {
		// name
		log.Println("name:" + s.ChildrenFiltered(`[property="v:itemreviewed"]`).Text())
		// year
		log.Println("year:" + s.ChildrenFiltered(`.year`).Text())
	})

	// #info > span:nth-child(1) > span.attrs
	director := ""
	doc.Find("#info span:nth-child(1) span.attrs").Each(func(i int, s *goquery.Selection) {
		// 导演
		director += s.Text()
	})
	log.Println("导演:" + director)

	pl := ""
	doc.Find("#info span:nth-child(3) span.attrs").Each(func(i int, s *goquery.Selection) {
		pl += s.Text()
	})
	log.Println("编剧:" + pl)

	charactor := ""
	doc.Find("#info span.actor span.attrs").Each(func(i int, s *goquery.Selection) {
		charactor += s.Text()
	})
	log.Println("主演:" + charactor)

	typeStr := ""
	doc.Find("#info > span:nth-child(8)").Each(func(i int, s *goquery.Selection) {
		typeStr += s.Text()
	})
	log.Println("类型:" + typeStr)

}

func GetToplist(url string) []string {
	var urls []string
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	//bodyString, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(bodyString))
	if resp.StatusCode != 200 {
		log.Println("err resp:", resp)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("#content div div.article ol li div div.info div.hd a").Each(func(i int, s *goquery.Selection) {
		// year
		log.Printf("%v", s)
		herf, _ := s.Attr("href")
		urls = append(urls, herf)
	})

	return urls
}

func StartDouban() {
	var urls []string
	var newUrl string

	url := "https://movie.douban.com/top250?start="
	log.Printf("%v", url)

	for i := 0; i < 10; i++ {
		start := i * 25
		newUrl = url + strconv.Itoa(start)
		urls = GetToplist(newUrl)

		for _, url := range urls {
			GetMovie(url)
		}
	}
}
