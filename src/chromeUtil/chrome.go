package chromeUtil

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
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
}

func NewTask(path string, url string, logInfo string, infoId int) *Task {
	return &Task{path: path, url: url, logInfo: logInfo, infoId: infoId}
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
		cs.startChromeService(ch)
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

	webDriver := getWebDriver()
	defer webDriver.Close()

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

	//判断加载完成
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

	log.Println("开始写文件:", task.path)
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

func getWebDriver() selenium.WebDriver {
	// 调起chrome浏览器
	webDriver, err := selenium.NewRemote(getChromeCaps(), fmt.Sprintf("http://localhost:%d/wd/hub", Port))
	if err != nil {
		log.Fatalln(err)
	}
	return webDriver
}

func getChromeCaps() selenium.Capabilities {
	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	// 禁止加载图片，加快渲染速度
	imgCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}

	chromeCaps := chrome.Capabilities{
		Prefs: imgCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
			"--start-maximized",
			"--disable-gpu",
			"--disable-impl-side-painting",
			"--disable-gpu-sandbox",
			"--disable-accelerated-2d-canvas",
			"--disable-accelerated-jpeg-decoding",
			"--test-type=ui",
		},
	}
	//以上是设置浏览器参数
	caps.AddChrome(chromeCaps)

	return caps
}
