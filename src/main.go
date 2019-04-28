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

	service := chromeUtil.NewChromeService()
	service.StartChromeService()
	log.Println("chromeService服务启动完毕!")

	defer service.Close()

	queue := util.NewSeqWaitModel(lessonsLen)
	defer queue.Close()

	for _, id := range allLessonIds {
		go func(id int) {
			infoTask := jike.NewInfoTask(id, queue.GetChan())
			queue.AddTask(infoTask)
			//jike.GetOneLessonInfo(id)
			//ch <- true
		}(id)
	}

	queue.Wait()
}
