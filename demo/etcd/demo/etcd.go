package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

// etcd client put/get demo
// use etcd/clientv3

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	fmt.Println("connect to etcd success")
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {
			panic(err)
		}
	}(cli)
	// put
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	value2 := `[{"log_path":"D:/tmp/log/nginx.log","topic":"web_log"},{"log_path":"D:/tmp/log/mysql.log","topic":"mysql_log"}]`
	value3 := `[{"log_path":"D:/tmp/log/nginx.log","topic":"web_log"},{"log_path":"D:/tmp/log/mysql.log","topic":"mysql_log"},{"log_path":"D:/tmp/log/redis.log","topic":"redis_log"}]`
	value := `[{"log_path":"D:/tmp/log/redis.log","topic":"redis_log"}]`
	_,_ = value,value2

	_, err = cli.Put(ctx, "/logagent/192.168.0.100/collect_config", value3)

	cancel()
	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
	//// get
	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//resp, err := cli.Get(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Printf("get from etcd failed, err:%v\n", err)
	//	return
	//}
	//for _, ev := range resp.Kvs {
	//	fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	//}
}
