package kafka

import (
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type l struct {
	topic   string
	message string
}

var (
	producer sarama.SyncProducer // 声明全局连接kafka的成产者)
	logChan  chan *l
)

//Init 初始化生产者(与init无关)
func Init(address []string, maxSize int) (err error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	producer, err = sarama.NewSyncProducer(address, config)   // 创建kafka消费者
	if err != nil {
		fmt.Println("create producer failed,err:", err)
		return
	}
	// 初始化 logChan
	logChan = make(chan *l, maxSize)
	// 开启后台goroutine从通道取数据发给kafka
	go sendToKafka()
	return
}

func SendToChan(topic string, msg string) {
	logChan <- &l{
		topic:   topic,
		message: msg,
	}
}

func sendToKafka() {
	for {
		select {
		case l := <-logChan:
			partition, offset, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: l.topic,
				Value: sarama.StringEncoder(l.message),
			})
			if err != nil {
				log.Fatalf("producer send message failed,err:%v\n", err)
			}
			fmt.Printf("patition:%d offset:%d\n", partition, offset)
		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}
