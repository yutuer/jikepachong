package jike

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Info struct {
	Cid           string `json:"cid"`
	With_groupbuy bool   `json:"with_groupbuy"`
}

func GetOneInfo(id int) {
	infoUrl := "https://time.geekbang.org/serv/v1/column/intro"

	info := &Info{Cid: strconv.Itoa(id), With_groupbuy: true}

	bs, err := json.Marshal(info)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(bs))

	request, e := http.NewRequest("POST", infoUrl, bytes.NewReader(bs))
	if e != nil {
		log.Fatalln(e)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", fmt.Sprintf("https://time.geekbang.org/column/intro/%d", id))
	request.Header.Set("Cookie", "_ga=GA1.2.1909908536.1541420176; GCID=137917f-117ceb1-f7292b8-7205755; _gid=GA1.2.1174986564.1555901483; _gat=1; Hm_lvt_022f847c4e3acd44d4a2481d9187f1e6=1555759694,1555901483,1555901486,1555937472; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1555937474|1555935529; Hm_lpvt_022f847c4e3acd44d4a2481d9187f1e6=1555937475")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	{
		log.Println(resp.StatusCode)
		if resp.StatusCode == 200 {
			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(string(bs))
		}

	}

}
