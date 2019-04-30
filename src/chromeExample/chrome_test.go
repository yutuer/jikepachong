package chromeExample

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"testing"
	"time"
)

func TestChrome(t *testing.T) {

	//如果seleniumServer没有启动，就启动一个seleniumServer所需要的参数，可以为空，示例请参见https://github.com/tebeka/selenium/blob/master/example_test.go
	opts := []selenium.ServiceOption{}
	//opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	//}

	//selenium.SetDebug(true)

	service, err := selenium.NewChromeDriverService(SeleniumPath, Port, opts...)
	if nil != err {
		t.Fatal("start a chromedriver service falid", err.Error())
	}
	//注意这里，server关闭之后，chrome窗口也会关闭
	defer service.Stop()

	//链接本地的浏览器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	//禁止图片加载，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 设置Chrome无头模式，在linux下运行，需要设置这个参数，否则会报错
			"--no-sandbox",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模拟user-agent，防反爬
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

	// 调起chrome浏览器
	w_b1, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Port))
	if err != nil {
		t.Fatal("connect to the webDriver faild", err.Error())
	}

	log.Println(w_b1)

	//关闭一个webDriver会对应关闭一个chrome窗口
	//但是不会导致seleniumServer关闭
	defer w_b1.Quit()

	err = w_b1.Get("https://time.geekbang.org/column/article/333")
	if err != nil {
		t.Fatal("get page faild", err.Error())
	}

	//// 重新调起chrome浏览器
	//w_b2, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", Port))
	//if err != nil {
	//	t.Fatal("connect to the webDriver faild", err.Error())
	//	return
	//}
	//defer w_b2.Close()
	////打开一个网页
	//err = w_b2.Get("https://www.toutiao.com/")
	//if err != nil {
	//	t.Fatal("get page faild", err.Error())
	//	return
	//}
	////打开一个网页
	//err = w_b2.Get("https://www.baidu.com/")
	//if err != nil {
	//	t.Fatal("get page faild", err.Error())
	//	return
	//}
	//w_b就是当前页面的对象，通过该对象可以操作当前页面了
	//........

	//判断加载完成

	jsRt, err := w_b1.ExecuteScript("return document.readyState", nil)
	if err != nil {
		log.Fatalln("exe js err", err)
	}
	if jsRt != "complete" {
		log.Fatalln("网页加载未完成")
	}

	log.Println("开始休眠2S...")

	time.Sleep(1 * time.Second)

	//data, err := w_b1.FindElement(selenium.ByXPATH, "//*[@id=\"app\"]/div[1]/div[2]/div[2]")
	data, err := w_b1.FindElement(selenium.ByID, "app")
	if err != nil {
		log.Fatalln(err)
	}

	s, err := data.Text()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(s)

	//bs, err := data.Screenshot(true)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//util.WriteFile_B(bs, "d:/111.png")

	//frameHtml, err := w_b1.PageSource()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(frameHtml)

	return
}
