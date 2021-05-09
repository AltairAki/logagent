package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

var client *clientv3.Client

// CollectEntry etcd中对应的配置项
type CollectEntry struct {
	LogPath string `json:"log_path"` // 日志path
	Topic   string `json:"topic"`    // 日志发往kafka的topic
}

func Init(address []string, timeout time.Duration) (err error) {
	config := clientv3.Config{
		Endpoints:   address,
		DialTimeout: timeout,
	}
	client, err = clientv3.New(config)
	return
}

// Get 根据key从etcd中获取配置项
func Get(key string) (collectEntryList []*CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	response, err := client.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatalf("etcd get key failed,err: %v", err)
	}
	for _, kv := range response.Kvs {
		fmt.Printf("%s:%s\n", kv.Key, kv.Value)
		err = json.Unmarshal(kv.Value, &collectEntryList)
		if err != nil {
			log.Fatalf("etcd config json unmarshal failed,err: %v", err)
		}
	}
	return collectEntryList, err
}

// Watch 如果etcd配置有更新就通知 taillog.manager.collectEntryChan
func Watch(key string, collectEntryChan chan<- []*CollectEntry) {
	watchChan := client.Watch(context.Background(), key)
	// 从通道取值
	for watch := range watchChan {
		for _, event := range watch.Events {
			fmt.Printf("Type:%s Key:%s Value:%s\n", event.Type, event.Kv.Key, event.Kv.Value)
			// 通知 taillog.manager
			var collectEntry []*CollectEntry
			// 如果是非删除则为collectEntry赋值
			if event.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(event.Kv.Value, &collectEntry)
				if err != nil {
					fmt.Printf("unmarshal failed,err:%v\n", err)
					continue
				}
			}
			fmt.Printf("etcd conf change and set collectEntry chan:%v\n", collectEntry)
			collectEntryChan <- collectEntry
		}
	}
}
