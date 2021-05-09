package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("子go程：执行第%d次操作！\n", i)
			time.Sleep(time.Second)
		}
	}()
	for i := 0; i < 3; i++ {
		fmt.Println("--------aaaa------")
		time.Sleep(time.Second)
	}

}
