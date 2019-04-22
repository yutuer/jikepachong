package jike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
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
	List ArticleList `json:"list"`
	Page ArticlePage `json:"page"`
}

type ArticleList struct {

}

type ArticlePage struct {
	Count int `json:"count"`
	More bool `json:"more"`
}

func GetArticle(id int) {
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

}
