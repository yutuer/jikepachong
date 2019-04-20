package example

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
)

func ExampleScrape1() {
	html := `<body>

				<div>DIV1</div>
				<div>DIV2</div>
				<span>SPAN</span>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("div").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})
}

func ExampleScrape2() {
	html := `<body>

				<div id="div1">DIV1</div>
				<div>DIV2</div>
				<span>SPAN</span>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("#div1").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})
}

func ExampleScrape3() {
	html := `<body>

				<div id="div1">DIV1</div>
				<div class="name">DIV2</div>
				<span>SPAN</span>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find(".name").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})
}

func ExampleScrape4() {
	html := `<body>

				<div>DIV1</div>
				<div class="name">DIV2</div>
				<div class="name1">DIV3</div>
				<span>SPAN</span>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	dom.Find("div[class]").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})

	dom.Find("div[class=name]").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})

	dom.Find("div[class=name1]").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})
}

func ExampleScrape5() {
	html := `<body>

				<div lang="ZH">DIV1</div>
				<div lang="zh-cn">DIV2</div>
				<div lang="en">DIV3</div>
				<span>
					<div>DIV4</div>
				</span>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	// 筛选出body这个父元素下，符合条件的最直接的子元素div，结果是DIV1、DIV2、DIV3，虽然DIV4也是body的子元素，但不是一级的，所以不会被筛选到。
	dom.Find("body>div").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})

	// 就是要筛选body下所有的div元素，不管是一级、二级还是N级。
	dom.Find("body div").Each(func(i int, selection *goquery.Selection) {
		log.Println(selection.Text())
	})

	//dom.Find("div").Each(func(i int, selection *goquery.Selection) {
	//	log.Println(selection.Text())
	//})
}

func ExampleScrape6() {
	html := `<body>

				<div lang="zh">DIV1</div>
				<p>P1</p>
				<div lang="zh-cn">DIV2</div>
				<div lang="en">DIV3</div>
				<span>
					<div>DIV4</div>
				</span>
				<p>P2</p>

			</body>
			`

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatalln(err)
	}

	//相邻的选择器
	//dom.Find("div[lang=zh]+p+div").Each(func(i int, selection *goquery.Selection) {
	//	fmt.Println(selection.Text())
	//})

	// 兄弟选择器
	//dom.Find("div[lang=zh]~p").Each(func(i int, selection *goquery.Selection) {
	//	fmt.Println(selection.Text())
	//})

	//内容筛选器
	//dom.Find("div:contains(DIV2)").Each(func(i int, selection *goquery.Selection) {
	//	fmt.Println(selection.Text())
	//})

	//Find(":has(selector)")和contains差不多，只不过这个是包含的是元素节点。
	dom.Find("span:has(div)").Each(func(i int, selection *goquery.Selection) {
		fmt.Println(selection.Text())
	})
}

func TestAll() {
	//ExampleScrape1()
	//ExampleScrape2()
	//ExampleScrape3()
	//ExampleScrape4()
	//ExampleScrape5()
	//ExampleScrape6()
}
