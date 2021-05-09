package main

import (
	"fmt"

	"github.com/Shopify/sarama"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	// 连接kafka
	producer, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed,err:", err)
		return
	}
	defer producer.Close()	 // 日志收集项目不需要关闭kafka连接

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("这是一条测试消息")
	// 发送消息
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("send message failed,err:", err)
		return
	}
	fmt.Printf("partition:%v offset:%v\n", partition, offset)
}
