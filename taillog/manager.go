package taillog

import (
	"fmt"
	"logagent/etcd"
)

var tailorManager *manager

// tailor 管理者
type manager struct {
	collectEntryList []*etcd.CollectEntry
	tailorMap        map[string]*Tailor
	collectEntryChan chan []*etcd.CollectEntry
}

// InitTask 为每个日志配置项都初始化一个tailor 并监听etcd是否有配置更新
// collectEntryList etcd 读到的配置条目
func InitTask(collectEntryList []*etcd.CollectEntry) {
	tailorManager = &manager{
		collectEntryList: collectEntryList,                // 保存当前的日志收集项目配置信息
		tailorMap:        make(map[string]*Tailor, 16),    // 最多存16种日志
		collectEntryChan: make(chan []*etcd.CollectEntry), // 无缓冲区通道；etcd配置没有修改就阻塞等着
	}
	// 将etcd中读到的每个日志配置都初始化一个 tailor
	for _, collectEntry := range collectEntryList {
		tailor := NewTailor(collectEntry.LogPath, collectEntry.Topic)
		mk := fmt.Sprintf("%s_%s", collectEntry.LogPath, collectEntry.Topic)
		tailorManager.tailorMap[mk] = tailor
		// 将文件读取的日志发到 kafka.logChan
		go tailor.sendToKafkaChan()
	}

	// 监听etcd配置的新增，修改，删除
	go tailorManager.listencollectEntryChan()
}

// 监听 collectEntryChan (监听etcd配置的新增，修改，删除)
func (m *manager) listencollectEntryChan() {
	for {
		select {
		case collectEntry := <-m.collectEntryChan:
			for _, entry := range collectEntry {
				mk := fmt.Sprintf("%s_%s", entry.LogPath, entry.Topic)
				// 如果tailorMap中没有对应的字段则说明没有对应的 Tailor
				if _, ok := m.tailorMap[mk]; ok {
					// 0. 原来就有，无需操作
					continue
				} else {
					// 1. 新增配置，初始化新的 Tailor 并更新 tailorMap
					tailor := NewTailor(entry.LogPath, entry.Topic)
					m.tailorMap[mk] = tailor
					go tailor.sendToKafkaChan()
				}
			}

			// 遍历原有配置与现有配置比较，把原有的但现在没有的从 m.tailorMap中删掉并取消掉
			for _, old := range m.collectEntryList {
				isFound := false
				for _, now := range collectEntry {
					if now.LogPath == old.LogPath && now.Topic == old.Topic {
						isFound = true
						continue
					}
				}
				// 该配置项在最新的配置中找不到，就应该取消
				mk := fmt.Sprintf("%s_%s", old.LogPath, old.Topic)
				if _, ok := m.tailorMap[mk]; ok && !isFound {
					m.tailorMap[mk].cancelFunc()
					delete(m.tailorMap, mk)
				}
			}
			fmt.Println("get a new etcd conf:", collectEntry)
		}
	}
}

func GetCollectEntryChan() chan<- []*etcd.CollectEntry {
	return tailorManager.collectEntryChan
}
