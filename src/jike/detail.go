package jike

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)


func StartDetail() {
	mainUrl := "https://time.geekbang.org/serv/v1/column/details"

	//tr := &http.Transport{
	//	TLSClientConfig:    &tls.Config{RootCAs: pool},
	//	DisableCompression: true,
	//}
	//client := &http.Client{Transport: tr}

	// 1m
	bs := make([]byte, 1*1024*1024)
	resp, err := http.Post(mainUrl, "application/json", bytes.NewReader(bs))
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(all)
}
