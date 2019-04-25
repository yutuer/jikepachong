package jike

import (
	"bytes"
	"chromeUtil"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
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
	return fmt.Sprintf("%d, %s, %d", article.Id, article.Article_title, article.Chapter_id)
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

			webDriver := getWebDriver()
			defer webDriver.Close()

			for _, v := range articleRes.Data.List {
				DoOneArticle(webDriver, dirPath, v, id)
			}
		} else {
			log.Println(resp.StatusCode)
		}
	}
}

func DoOneArticle(driver selenium.WebDriver, dirPath string, article Article, infoId int) {
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

	task := chromeUtil.NewTask(path, url, article.String(), infoId, driver)

	chromeUtil.GetChromeService().SubmitTask(task)
}

func getWebDriver() selenium.WebDriver {
	// 调起chrome浏览器
	webDriver, err := selenium.NewRemote(getChromeCaps(), fmt.Sprintf("http://localhost:%d/wd/hub", Port))
	if err != nil {
		log.Fatalln(err)
	}

	return webDriver
}

func getChromeCaps() selenium.Capabilities {
	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imgCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			//"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
			//"--start-maximized",
			//"--disable-gpu",
			//"--disable-impl-side-painting",
			//"--disable-gpu-sandbox",
			//"--disable-accelerated-2d-canvas",
			//"--disable-accelerated-jpeg-decoding",
			//"--test-type=ui",
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	return caps
}
