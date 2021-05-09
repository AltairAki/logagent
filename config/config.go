package config

import "github.com/go-ini/ini"

type Note struct {
	Kafka `ini:"kafka"`
	Etcd  `ini:"etcd"`
}

type Kafka struct {
	Address     []string `ini:"address"`
	ChanMaxSize int      `ini:"chan_max_size"`
}

type Etcd struct {
	Address []string `ini:"address"`
	CollectLogKey     string   `ini:"collect_log_key"`
	Timeout int      `ini:"timeout"`
}

// Taillog -------unused---------
// filename 从 etcd 中获取
type Taillog struct {
	Filename string `ini:"filename"`
}

var Cfg = new(Note) // 注意这里必须是new不然仅声明Cfg是nil，会报错：reflect: call of reflect.Value.Type on zero Value

func NewConfig(path string) (cfg *Note, err error) {
	err = ini.MapTo(Cfg, path)
	return Cfg, err
}
