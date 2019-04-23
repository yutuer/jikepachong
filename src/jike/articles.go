package jike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	FileName   = "%s/%s.html"
	ArticleUrl = "https://time.geekbang.org/column/article/%d"
)

type ArticleReq struct {
	Cid    string `json:"cid"`
	Order  string `json:"order"`
	Prev   int    `json:"prev"`
	Sample bool   `json:"sample"`
	Size   int    `json:"size"`
}

type ArticleRes struct {
	Code int         `json:"code"`
	Data ArticleData `json:"data"`
}

type ArticleData struct {
	List []Article   `json:"list"`
	Page ArticlePage `json:"page"`
}

type Article struct {
	Id            int    `json:"id"`
	Article_title string `json:"article_title"`
	Chapter_id    int    `json:"chapter_id"`
}

type ArticlePage struct {
	Count int  `json:"count"`
	More  bool `json:"more"`
}

func GetArticles(dirPath string, id int) {
	articleUrl := "https://time.geekbang.org/serv/v1/column/articles"

	req := &ArticleReq{Cid: strconv.Itoa(id), Order: "earliest", Prev: 0, Sample: false, Size: 500}

	bs, e := json.Marshal(req)
	if e != nil {
		log.Fatalln(e)
	}

	request, e := http.NewRequest("POST", articleUrl, bytes.NewReader(bs))
	if e != nil {
		log.Fatalln(e)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", fmt.Sprintf("https://time.geekbang.org/column/intro/%d", id))

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	{
		if resp.StatusCode == 200 {
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}

			articleRes := &ArticleRes{}
			json.Unmarshal(bs, articleRes)

			for _, v := range articleRes.Data.List {
				DoOneArticle(dirPath, v, id)
			}
		} else {
			log.Println(resp.StatusCode)
		}
	}
}

func DoOneArticle(dirPath string, article Article, infoId int) {
	url := fmt.Sprintf(ArticleUrl, article.Id)
	title := strings.Replace(article.Article_title, "/", "&", -1)
	title = strings.Replace(title, "27 | ", "27__", -1)
	title = strings.Replace(title, " ", "", -1)
	title = strings.Replace(title, "？", "!", -1)
	title = strings.Replace(title, "“", "", -1)
	title = strings.Replace(title, "”", "", -1)
	title = strings.Replace(title, "*", "_", -1)
	title = strings.Replace(title, "?", "!", -1)
	title = strings.Replace(title, "：", "_", -1)
	title = strings.Replace(title, "-", "_", -1)
	title = strings.Replace(title, "|", "__", -1)
	path := fmt.Sprintf(FileName, dirPath, title)
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if e != nil {
		log.Fatalln(e)
	}

	defer f.Close()

	log.Printf("开始抓取 %d: %v", infoId, article)

	content := StartChromeAndGetContent(url)
	f.Write([]byte(content))

	// 暂停2秒
	//time.Sleep(2 * time.Second)
}
