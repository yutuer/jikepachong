package chrome

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"time"
)

const (
	port = 9515
)

// StartChrome 启动谷歌浏览器headless模式
func StartChromeAndGetContent(url string) string {
	//如果seleniumServer没有启动，就启动一个seleniumServer所需要的参数，可以为空，示例请参见https://github.com/tebeka/selenium/blob/master/example_test.go
	opts := []selenium.ServiceOption{
	//opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	}

	//selenium.SetDebug(true)

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
			"--start-maximized",
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_2) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0.2 Safari/604.4.7", // 模拟user-agent，防反爬
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

	// 启动chromedriver，端口号可自定义
	service, err := selenium.NewChromeDriverService("chromedriver.exe", port, opts...)
	if err != nil {
		log.Fatalln("Error starting the ChromeDriver server: %v", err)
	}

	//注意这里，server关闭之后，chrome窗口也会关闭
	defer service.Stop()

	// 调起chrome浏览器
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Fatalln(err)
	}

	//关闭一个webDriver会对应关闭一个chrome窗口
	//但是不会导致seleniumServer关闭
	defer webDriver.Close()

	// 这是目标网站留下的坑，不加这个在linux系统中会显示手机网页，每个网站的策略不一样，需要区别处理。
	webDriver.AddCookie(&selenium.Cookie{
		Name:  "defaultJumpDomain",
		Value: "www",
	})

	// 导航到目标网站
	err = webDriver.Get(url)
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

	//获取网站内容
	time.Sleep(1 * time.Second)
	frameHtml, err := webDriver.PageSource()
	if err != nil {
		log.Fatalln(err)
	}

	return frameHtml
}
