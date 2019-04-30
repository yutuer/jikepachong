package chromedpExample

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"testing"
)

const (
	Url = "https://time.geekbang.org/column/article/12655"
	//Url = "https://juejin.im/post/5cbd79e85188250a7f630baa"
)

func TestA(t *testing.T) {
	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://golang.org/pkg/time/`),
		chromedp.Text(`#pkg-overview`, &res, chromedp.NodeVisible, chromedp.ByID),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(strings.TrimSpace(res))
}

func TestB(t *testing.T) {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		//chromedp.Headless,
		chromedp.Flag("headless", false),
		chromedp.DisableGPU,
		//chromedp.UserDataDir(dir),
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	contentMustTemplate := `//*[@id="app"]/div[1]/div[2]/div[2]/div[1]/div[2]`
	clearTemplate := `//*[@id="app"]//div[@class="_3-b6SqNP_0"]/*`
	//contentTemplate := `//*[@id="app"]`
	//contentTemplate := `/html`
	//cssTemplate := ``

	// ensure that the browser process is started
	//var css string
	var content string

	err := chromedp.Run(taskCtx,
		chromedp.Navigate(Url),
		chromedp.WaitVisible(contentMustTemplate, chromedp.BySearch),
		chromedp.Evaluate(clearTemplate, chromedp.BySearch),
		//chromedp.OuterHTML(contentTemplate, &content, chromedp.BySearch),
		//chromedp.OuterHTML(cssTemplate, &css, chromedp.BySearch),
	)

	fmt.Println(content)

	if err != nil {
		log.Fatalln(err)
	}

}
