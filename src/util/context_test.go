package util

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestChannel(t *testing.T) {
	stop := make(chan struct{})

	go watchChannel(stop, "【监控1】")
	go watchChannel(stop, "【监控2】")
	go watchChannel(stop, "【监控3】")

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	stop <- struct{}{}
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func watchChannel(stop chan struct{}, name string) {
	for {
		select {
		case x := <-stop: //用channel 只会有一个读取, 其他都读取不到
			fmt.Println(name, "监控退出，停止了...  <-done is:", x)
			return
		default:
			fmt.Println(name, "goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}

func TestContext1(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func TestContext2(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go watchContext(ctx, "【监控1】")
	go watchContext(ctx, "【监控2】")
	go watchContext(ctx, "【监控3】")

	fmt.Println(ctx.Deadline())

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()

	fmt.Println(ctx.Err())
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func watchContext(ctx context.Context, name string) {
	for {
		select {
		case x := <-ctx.Done(): //用context 全部都可以读取
			fmt.Println(name, "监控退出，停止了...  <-done is:", x)
			return
		default:
			fmt.Println(name, "goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}
