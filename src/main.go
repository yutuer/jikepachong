package main

import (
	"jike"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	//GCID, GCESS, SERVERID := jike.StartAccount()
	//log.Println(GCID, GCESS, SERVERID)

	allIds := jike.StartNewAll()
	log.Println("num:", len(allIds), ", allIds:", allIds)

	ch := make(chan bool, len(allIds))

	for _, id := range allIds {
		jike.GetOneInfo(id)
		ch <- true
	}

	for i := 0; i < len(ch); i++ {
		<-ch
	}
}
