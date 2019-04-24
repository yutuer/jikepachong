package jike

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type List struct {
	//Column_ctime        int  `json:"column_ctime"`
	//Column_groupbuy     int  `json:"column_groupbuy"`
	//Column_price        int  `json:"column_price"`
	//Column_price_market int  `json:"column_price_market"`
	//Column_sku          int  `json:"column_sku"`
	//Column_type         int  `json:"column_type"`
	//Had_sub             bool `json:"had_sub"`
	Id int `json:"id"`
	//Is_experience       bool `json:"is_experience"`
	//Last_aid            int  `json:"last_aid"`
	//Price_type          int  `json:"price_type"`
	//Sub_count           int  `json:"sub_count"`
}

type Nav struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
	Id    int    `json:"id"`
	Name  string `json:"name"`
}

type Page struct {
	Count int `json:"count"`
}

type Data struct {
	List []List `json:"list"`
	Nav  []Nav  `json:"nav"`
	Page *Page  `json:"page"`
}

type NewAllResult struct {
	Error []string `json:"error"`
	Extra []string `json:"extra"`
	Data  Data     `json:"data"`
	Code  int      `json:"code"`
}

func GetLessonIds() []int {
	mainUrl := "https://time.geekbang.org/serv/v1/column/newAll"

	client := &http.Client{}

	request, err := http.NewRequest("POST", mainUrl, strings.NewReader(""))
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", "https://time.geekbang.org/?category=0") //重要
	//request.Header.Set("Cookie", "_ga=GA1.2.1909908536.1541420176; GCID=137917f-117ceb1-f7292b8-7205755; _gid=GA1.2.1189370068.1555658102; _gat=1; GCESS=BAUEAAAAAAQEAC8NAAcEVq8nPgsCBAAIAQMBBFUPEwADBP7LuVwMAQEGBPzgQgQKBAAAAAACBP7LuVwJAQE-; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1555680254|1555675541") //重要
	//request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	//request.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	ret := make([]int, 0)

	if resp.StatusCode == 200 {
		// 生成对象
		NewAllResult := &NewAllResult{}

		json.Unmarshal(bs, NewAllResult)

		for _, v := range NewAllResult.Data.List {
			ret = append(ret, v.Id)
		}
	} else {
		log.Println(resp.StatusCode)
	}
	return ret
}
