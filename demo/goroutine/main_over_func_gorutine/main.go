package main

import (
	"fmt"
	"runtime"
	"time"
)

func test() {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("test 执行第%d次goroutine\n", i)
			time.Sleep(time.Second)
		}
	}()
	fmt.Println("test goroutine执行完毕！")
}

func main() {
	test()
	for {
		runtime.GC()
	}
}
