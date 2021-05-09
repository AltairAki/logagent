package main

import (
	"fmt"
	"runtime"
	"time"
)

func test() {
	for i := 0; i < 10; i++ {
		fmt.Printf("test 执行第%d次goroutine\n", i)
		time.Sleep(time.Second)
	}
	fmt.Println("test goroutine执行完毕！")
}
func main() {
	go func() {
		go test()
		fmt.Println("------aaaaaaaa-------")
		time.Sleep(time.Second)
		//fmt.Println("------goroutine结束--------------")
		/*
		   不管是return  还是  runtime.Goexit()，效果一样
		*/
		return
		//runtime.Goexit()
	}()
	for {
		runtime.GC()
	}
}
