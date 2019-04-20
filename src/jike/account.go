package jike

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Account struct {
	Appid     int    `json:"appid"`
	Captcha   string `json:"captcha"`
	Cellphone string `json:"cellphone"`
	Country   int    `json:"country"`
	Password  string `json:"password"`
	Platform  int    `json:"platform"`
	Remember  int    `json:"remember"`
}

func StartAccount() (string, string, string) {
	account := &Account{
		Appid:     1,
		Captcha:   "",
		Cellphone: "13699105433",
		Country:   86,
		Password:  "1980514zf",
		Platform:  3,
		Remember:  1,
	}

	bs, err := json.MarshalIndent(account, "", "    ")
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(bs))

	accountUrl := "https://account.geekbang.org/account/ticket/login"

	request, err := http.NewRequest("POST", accountUrl, bytes.NewReader(bs))
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Referer", "https://account.geekbang.org/login") //必须填写
	//request.Header.Set("Cookie", "_ga=GA1.2.1909908536.1541420176; GCID=137917f-117ceb1-f7292b8-7205755; _gid=GA1.2.1189370068.1555658102; SERVERID=3431a294a18c59fc8f5805662e2bd51e|1555751107|1555750718")
	//request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	{
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
	}

	var GCID string
	var GCESS string
	var SERVERID string

	for _, v := range resp.Cookies() {
		if v.Name == "GCID" {
			GCID = v.Value
		} else if v.Name == "GCESS" {
			GCESS = v.Value
		} else if v.Name == "SERVERID" {
			SERVERID = v.Value
		}
	}

	return GCID, GCESS, SERVERID
}
