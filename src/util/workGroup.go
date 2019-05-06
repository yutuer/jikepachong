package util

import (
	"fmt"
	"sync"
)

func exam1() {
	total := 100

	wg := sync.WaitGroup{}
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func exam2() {
	total := 100
	diff := -10

	wg := sync.WaitGroup{}
	//wg 添加一个负数
	//1. 先添加100, 在加-10是正确的, Add参数可以为负数, 运行90次
	wg.Add(total)
	wg.Add(diff)

	//2. 先添加-10, 在加100是错误的
	//wg.Add(diff)
	//wg.Add(total)

	for i := 0; i < total+diff; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

// 可以重复利用
func exam3() {
	total := 100

	wg := sync.WaitGroup{}
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println("once again!!")

	wg.Add(total)
	// 周期走完, 再次利用
	for i := 0; i < total; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

}

// 利用future来获取结果
func exam4() {
	const total = 10

	var results [total]int

	wg := sync.WaitGroup{}
	wg.Add(total)

	for i := 0; i < total; i++ {
		go func(i int) {
			fmt.Println(i)
			wg.Done()

			results[i] = i
		}(i)
	}

	wg.Wait()

	fmt.Println(results)
}

func exam5() {
	once := sync.Once{}

	go once.Do(func() { fmt.Println(1) })
	go once.Do(func() { fmt.Println(2) })

}
