package jike

import (
	"bytes"
	"chromeUtil"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func (article *Article) String() string {
	return fmt.Sprintf("%d, %s", article.Id, article.Article_title)
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

			doArticles(articleRes, dirPath, id)
		} else {
			log.Println(resp.StatusCode)
		}
	}
}

func doArticles(articleRes *ArticleRes, dirPath string, id int) {
	length := len(articleRes.Data.List)

	ch := make(chan bool, length)


	for _, v := range articleRes.Data.List {
		func(article Article) {
			DoOneArticle_SendToQueue(dirPath, article, id)
			ch <- true
		}(v)
	}

	for i := 0; i < length; i++ {
		<-ch
	}
}

func DoOneArticle_SendToQueue(dirPath string, article Article, infoId int) {
	url := fmt.Sprintf(ArticleUrl, article.Id)
	title := strings.Replace(article.Article_title, "/", "&", -1)
	//title = strings.Replace(title, "27 | ", "27__", -1)
	title = strings.Replace(title, " ", "", -1)
	//title = strings.Replace(title, "？", "!", -1)
	//title = strings.Replace(title, "“", "", -1)
	//title = strings.Replace(title, "”", "", -1)
	//title = strings.Replace(title, "*", "_", -1)
	//title = strings.Replace(title, "?", "!", -1)
	//title = strings.Replace(title, "：", "_", -1)
	//title = strings.Replace(title, "-", "_", -1)
	title = strings.Replace(title, "|", "__", -1)

	path := fmt.Sprintf(FileName, dirPath, title)

	task := chromeUtil.NewTask(path, url, article.String(), infoId)

	chromeUtil.GetChromeService().SubmitTask(task)
}