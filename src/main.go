package main

import (
	"chromeUtil"
	"jike"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	//GCID, GCESS, SERVERID := jike.StartAccount()
	//log.Println(GCID, GCESS, SERVERID)

	allLessonIds := jike.GetLessonIds()
	lessonsLen := len(allLessonIds)
	log.Println("num:", lessonsLen, ", allIds:", allLessonIds)

	service := chromeUtil.NewChromeService()
	service.StartChromeService()
	log.Println("chromeService服务启动完毕!")

	defer service.Close()

	ch := make(chan bool, lessonsLen)
	for _, id := range allLessonIds {
		go func(id int) {
			jike.GetOneLessonInfo(id)
			ch <- true
		}(id)
	}

	for i := 0; i < lessonsLen; i++ {
		<-ch
	}

}
