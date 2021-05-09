package taillog

import (
	"logagent/kafka"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	log     *logrus.Logger
	localIP string
)

type Tailor struct {
	logPath    string
	topic      string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func init() {
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
}

// NewTailor  Tailor 构造函数
func NewTailor(logPath string, topic string) (t *Tailor) {
	// 初始化 Tailor.logPath Tailor.topic
	t = &Tailor{
		logPath: logPath,
		topic:   topic,
	}
	t.ctx, t.cancelFunc = context.WithCancel(context.Background())

	// 初始化tailFile
	// 每有一个新的日志都应该初始化一个新的 Tailor 负责追踪日志文件内容
	config := tail.Config{
		ReOpen:    true,                                 // 重新打开
		Follow:    true,                                 // 是否跟随
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // 读位置定位(重新打开文件也有定位)
		MustExist: false,                                // 日志文件是否必须存在，不存在会报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.logPath, config)
	if err != nil {
		fmt.Printf("tail file failed,err:%v\n", err)
		return
	}
	log.Infof("%s_%s tailor init success.", t.logPath, t.topic)
	return t
}

// 将文件读取的日志发到 kafka.logChan
func (t *Tailor) sendToKafkaChan() {
	for {
		select {
		case line := <-t.instance.Lines:
			kafka.SendToChan(t.topic, line.Text)
		case <-t.ctx.Done():
			log.Warnf("the task for path:%s is stop...", t.logPath)
			return // 函数返回对应的goroutine就结束了
		}
	}
}
