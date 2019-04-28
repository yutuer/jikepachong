package chromeUtil

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
)

const (
	Port         = 9515
	SeleniumPath = "chromedriver.exe"
)

type ChromeService struct {
	service *selenium.Service
}

func NewChromeService() *ChromeService {
	return &ChromeService{}
}

//注意 需要关闭返回的service
func (cs *ChromeService) StartChromeService() {
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
}

func (cs *ChromeService) Close() {
	if cs == nil {
		return
	}

	err := cs.service.Stop()
	if err != nil {
		log.Fatalln(err)
	}
}

func GetWebDriver() selenium.WebDriver {
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
