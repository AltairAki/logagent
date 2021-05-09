package main

import (
	"context"
	"fmt"
	"time"
)

type numbers struct {
	Num int
}

func counter(ctx context.Context) <-chan *numbers {
	out := make(chan *numbers)
	n := new(numbers)
	i := 0
	for ; ; i++ {
		n.Num = i
		out <- n
		if i == 10 {
			break
		}
	}
	close(out)
	return out
}

func main() {
	defer func() {
		fmt.Println("he")
	}()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	out := counter(ctx)
	for num := range out {
		fmt.Println(num)
	}
}
