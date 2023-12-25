package spider

import (
	"go-spider/config"
	"go-spider/src/utils/spider/douban"
	"go-spider/src/utils/spider/tiankong"
	"go-spider/src/utils/spider/toplist"
)

type SpiderTask interface {
	Start()
	PageDetail(id string)
	DoRecentUpdate()
}

// 定义 mod 的映射关系
var spiderModMap = map[string]SpiderTask{
	"async": &tiankong.AsyncSpiderApi{}, // 异步 goroutine 并行
	"sync":  &tiankong.SyncSpiderApi{},  // 同步 按顺序执行 串行
}

func Tiankong() SpiderTask {
	mod := config.ConfData.App.SpiderMod
	return spiderModMap[mod]
}

func Douban() SpiderTask {
	return &douban.DoubanApi{}
}

func Toplist() SpiderTask {
	return &toplist.ToplistApi{}
}
