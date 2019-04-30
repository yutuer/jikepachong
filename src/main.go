package main

import (
	"chromeUtil"
	"jike"
	"log"
	"util"
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

	service := chromeExample.NewChromeService()
	service.StartChromeService()
	log.Println("chromeService服务启动完毕!")

	defer service.Close()

	queue := util.NewNoSeqWaitModel(lessonsLen)
	defer queue.Close()

	for _, id := range allLessonIds {
		infoTask := jike.NewInfoTask(id)
		queue.AddTask(infoTask)
	}

	queue.Wait()
}
