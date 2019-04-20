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

	for _, id := range allIds {

	}
}
