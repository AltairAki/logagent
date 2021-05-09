package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		panic(err)
		return
	}
	fmt.Println("connect to etcd success")
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			panic(err)
		}
	}(cli)
	// watch key:name change(create，update，delete)
	watchChan := cli.Watch(context.Background(), "name")
	// try gey value from watchChan
	for watch := range watchChan {
		for _, ev := range watch.Events {
			fmt.Printf("Type: %s Key:%s Value:%s\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}
