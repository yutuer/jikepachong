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
)

const (
	Port   = 9515
	OutDir = "G:/jike/%d_%s_%s"
)

type InfoReq struct {
	Cid           string `json:"cid"`
	With_groupbuy bool   `json:"with_groupbuy"`
}

type InfoRes struct {
	Error []string `json:"error"`
	Extra []string `json:"extra"`
	Data  InfoData `json:"data"`
	Code  int      `json:"code"`
}

type InfoData struct {
	Author_name string `json:"author_name"`
	//Column_share_title string `json:"column_share_title"`
	//Author_intro       string `json:"author_intro"`
	//Column_subtitle    string `json:"column_subtitle"`
	Column_title string `json:"column_title"`
	//Column_unit        string `json:"column_unit"`
}

func GetOneLessonInfo(id int) {
	infoUrl := "https://time.geekbang.org/serv/v1/column/intro"

	info := &InfoReq{Cid: strconv.Itoa(id), With_groupbuy: true}

	bs, err := json.Marshal(info)
	if err != nil {
		log.Fatalln(err)
	}

	request, e := http.NewRequest("POST", infoUrl, bytes.NewReader(bs))
	if e != nil {
		log.Fatalln(e)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", fmt.Sprintf("https://time.geekbang.org/column/intro/%d", id))
	//request.Header.Set("Cookie", "_ga=GA1.2.1909908536.1541420176; GCID=137917f-117ceb1-f7292b8-7205755; _gid=GA1.2.1174986564.1555901483; _gat=1; Hm_lvt_022f847c4e3acd44d4a2481d9187f1e6=1555759694,1555901483,1555901486,1555937472; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1555937474|1555935529; Hm_lpvt_022f847c4e3acd44d4a2481d9187f1e6=1555937475")
	//request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

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

			infoRes := &InfoRes{}
			json.Unmarshal(bs, infoRes)

			dirPath := fmt.Sprintf(OutDir, id, infoRes.Data.Author_name, infoRes.Data.Column_title)

			if exists := isDirExist(dirPath); !exists {
				err = os.MkdirAll(dirPath, os.ModePerm)
				if err != nil {
					log.Fatalln(err)
				}
			}

			GetArticles(dirPath, id)
		} else {
			log.Println(resp.StatusCode)
		}

	}
}

func isDirExist(dirPath string) bool {
	_, e := os.Stat(dirPath)
	if e != nil {
		if os.IsExist(e) {
			return true
		} else {
			return false
		}
	}
	return true
}
