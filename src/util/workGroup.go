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
