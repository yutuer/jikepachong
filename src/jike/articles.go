package jike

import (
	"bytes"
	"chromeUtil"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"util"
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

func GetArticles(dirPath string, infoId int) error {
	articleUrl := "https://time.geekbang.org/serv/v1/column/articles"

	req := &ArticleReq{Cid: strconv.Itoa(infoId), Order: "earliest", Prev: 0, Sample: false, Size: 500}

	bs, e := json.Marshal(req)
	if e != nil {
		log.Fatalln(e)
	}

	request, e := http.NewRequest("POST", articleUrl, bytes.NewReader(bs))
	if e != nil {
		log.Fatalln(e)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", fmt.Sprintf("https://time.geekbang.org/column/intro/%d", infoId))

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

			doArticles(articleRes, dirPath, infoId)
		} else {
			log.Println(resp.StatusCode)
		}
	}

	return nil
}

func doArticles(articleRes *ArticleRes, dirPath string, infoId int) {
	length := len(articleRes.Data.List)

	driver := chromeUtil.GetWebDriver()
	defer driver.Close()

	queue := util.NewNoSeqWaitModel(length)
	defer queue.Close()

	for _, article := range articleRes.Data.List {
		//DoOneArticle_SendToQueue(dirPath, v, id)
		url := fmt.Sprintf(ArticleUrl, article.Id)
		title := strings.Replace(article.Article_title, "/", "&", -1)
		title = strings.Replace(title, " ", "", -1)
		title = strings.Replace(title, "|", "__", -1)
		path := fmt.Sprintf(FileName, dirPath, title)

		task := newTask(path, url, article.String(), infoId, driver)
		queue.AddTask(task)
	}

	queue.Wait()
}

func DoOneArticle_SendToQueue(dirPath string, article Article, infoId int) {
	//url := fmt.Sprintf(ArticleUrl, article.Id)
	//title := strings.Replace(article.Article_title, "/", "&", -1)
	////title = strings.Replace(title, "27 | ", "27__", -1)
	//title = strings.Replace(title, " ", "", -1)
	////title = strings.Replace(title, "？", "!", -1)
	////title = strings.Replace(title, "“", "", -1)
	////title = strings.Replace(title, "”", "", -1)
	////title = strings.Replace(title, "*", "_", -1)
	////title = strings.Replace(title, "?", "!", -1)
	////title = strings.Replace(title, "：", "_", -1)
	////title = strings.Replace(title, "-", "_", -1)
	//title = strings.Replace(title, "|", "__", -1)
	//
	//path := fmt.Sprintf(FileName, dirPath, title)
	//newTask(path, url, article.String(), infoId)
}

func newTask(path string, url string, logInfo string, infoId int, driver selenium.WebDriver) util.ITask {
	return &ArticleTask{path: path, url: url, logInfo: logInfo, infoId: infoId, webDriver: driver}
}

type ArticleTask struct {
	path      string
	url       string
	logInfo   string
	infoId    int
	webDriver selenium.WebDriver
}

// 启动一个新进程打开url,  返回内容
func (task *ArticleTask) DoTask() (error) {
	webDriver := task.webDriver

	log.Printf("开始抓取 %d: %s, %s", task.infoId, task.url, task.logInfo)
	// 导航到目标网站
	err := webDriver.Get(task.url)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to load page: %s\n", err))
	}

	//判断加载完成
	jsRt, err := webDriver.ExecuteScript("return document.readyState", nil)
	if err != nil {
		log.Fatalln("exe js err", err)
	}
	if jsRt != "complete" {
		log.Fatalln("网页加载未完成")
	}

	time.Sleep(1 * time.Second)

	//获取网站内容
	frameHtml, err := webDriver.PageSource()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("开始写文件:", task.path)
	util.WriteFile(frameHtml, task.path)

	return nil
}
