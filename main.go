package main

import (
	"fmt"
	"logagent/config"
	"logagent/etcd"
	"logagent/kafka"
	"logagent/taillog"
	"logagent/utils"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	log *logrus.Logger
)

func initLogger() {
	log = logrus.New()
	// 设置日志输出为os.Stdout
	log.Out = os.Stdout
	log.Level = logrus.DebugLevel
	// 可以设置像文件等任意`io.Writer`类型作为日志输出
	// file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	//  log.Out = file
	// } else {
	//  log.Info("Failed to log to file, using default stderr")
	// }
	log.Info("init log success.")
}

func main() {
	initLogger()
	// 0. 加载配置
	cfg, err := config.NewConfig("./config/config.ini")
	if err != nil {
		fmt.Println("config load failed,err:", err)
	}
	log.Info("loading config.ini success.")

	// 1. 初始化kafka producer
	err = kafka.Init(cfg.Kafka.Address, cfg.Kafka.ChanMaxSize)
	if err != nil {
		panic(err)
	}
	log.Info("init kafka success.")
	// 2. 初始化etcd
	err = etcd.Init(cfg.Etcd.Address, time.Duration(cfg.Etcd.Timeout)*time.Second)
	if err != nil {
		panic(err)
	}
	log.Info("init etcd success.")
	// 2.1 从etcd获取日志收集项目的配置信息
	ip, err := utils.GetOutboundIP()
	if err != nil {
		panic(fmt.Sprintf("get local ip failed, err:%v", err))
	}
	// 根据本机IP获取要收集日志的配置
	collectLogKey := fmt.Sprintf(cfg.Etcd.CollectLogKey, ip)
	log.Infof("collect log key:%s", collectLogKey)
	collectEntryList, err := etcd.Get(collectLogKey)
	log.Info("get conf from etcd success.")

	// 3. 收集日志并发给kafka
	taillog.InitTask(collectEntryList)

	var wg sync.WaitGroup
	wg.Add(1)
	// 2.2 开一个后台goroutine哨兵监控etcd日志收集配置变化(有变化及时通知logAgent实现热加载配置项)
	go etcd.Watch(collectLogKey, taillog.GetCollectEntryChan())
	wg.Wait()
}
