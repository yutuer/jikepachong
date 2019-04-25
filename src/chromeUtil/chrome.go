package chromeUtil

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"log"
	"os"
	"time"
)

const (
	Port         = 9515
	SeleniumPath = "chromedriver.exe"
)

type ChromeService struct {
	tasks   chan *Task
	service *selenium.Service
}

type Task struct {
	path    string
	url     string
	logInfo string
	infoId  int
	driver  selenium.WebDriver
}

func NewTask(path string, url string, logInfo string, infoId int, driver selenium.WebDriver) *Task {
	return &Task{path: path, url: url, logInfo: logInfo, infoId: infoId, driver: driver}
}

var service *ChromeService

func NewChromeService() *ChromeService {
	if service == nil {
		service = &ChromeService{tasks: make(chan *Task, 10000)}
	}
	return service
}

func GetChromeService() *ChromeService {
	return service
}

func (cs *ChromeService) StartChromeService() {
	ch := make(chan bool)

	go func() {
		go cs.startChromeService(ch)
	}()

	<-ch
}

//注意 需要关闭返回的service
func (cs *ChromeService) startChromeService(ch chan bool) {
	//如果seleniumServer没有启动，就启动一个seleniumServer所需要的参数，可以为空，示例请参见https://github.com/tebeka/selenium/blob/master/example_test.go
	opts := []selenium.ServiceOption{
		//opts := []selenium.ServiceOption{
		//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
		//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	}

	//selenium.SetDebug(true)

	// 启动chromedriver，端口号可自定义
	service, err := selenium.NewChromeDriverService(SeleniumPath, Port, opts...)
	if err != nil {
		log.Fatalf("Error starting the ChromeDriver server: %v", err)
	}

	cs.service = service

	ch <- true

	for task := range cs.tasks {
		executeTask(task)
	}
}

// 启动一个新进程打开url,  返回内容
func executeTask(task *Task) string {
	webDriver := task.driver

	log.Printf("开始抓取 %d: %s, %s", task.infoId, task.logInfo, task.url)
	// 导航到目标网站
	err := webDriver.Get(task.url)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to load page: %s\n", err))
	}

	id := webDriver.SessionID()
	log.Println("sessionId:", id)

	////判断加载完成
	jsRt, err := webDriver.ExecuteScript("return document.readyState", nil)
	if err != nil {
		log.Fatalln("exe js err", err)
	}
	if jsRt != "complete" {
		log.Fatalln("网页加载未完成")
	}

	time.Sleep(1 * time.Second)

	//获取网站内容
	frameHtml, err := webDriver.PageSource()
	if err != nil {
		log.Fatalln(err)
	}

	WriteFile(frameHtml, task.path)

	return frameHtml
}

func (cs *ChromeService) SubmitTask(task *Task) {
	cs.tasks <- task
}

func (cs *ChromeService) Close() {
	if cs == nil {
		return
	}

	close(cs.tasks)

	err := cs.service.Stop()
	if err != nil {
		log.Fatalln(err)
	}
}

func WriteFile(content string, path string) {
	var f *os.File

	_, err := os.Stat(path)
	if err == nil {
		err = os.Remove(path)
		if err != nil {
			log.Fatalln(err)
		}

		f, err = os.Create(path)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if os.IsNotExist(err) {
			f, err = os.Create(path)
			if err != nil {
				log.Fatalln(err)
			}
		}
	}

	if f == nil {
		log.Fatalln(errors.Errorf("title:%s 不能生成文件", path))
	}
	defer f.Close()

	f.Write([]byte(content))
}
