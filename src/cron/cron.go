package cron

import (
	"context"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go-spider/package/storage"
	"go-spider/src/utils/spider"
	"log"
)

func InitCron() {
	firstSpider()
	// 启动定时爬虫任务 全量
	TimingSpider(func() {
		spider.Tiankong().Start()
		return
	})

	// 爬虫 只爬取最近有更新的资源
	RecentUpdate(func() {
		spider.Tiankong().DoRecentUpdate()
		return
	})

	// 豆瓣启动定时爬虫任务 全量
	TimingSpider(func() {
		spider.Douban().Start()
		return
	})

	// 爬虫 只爬取最近有更新的资源
	TimingSpider(func() {
		spider.Toplist().Start()
		return
	})
}

// 首次启动自动开启爬虫
func firstSpider() {
	log.Println("hasHome1")
	hasHome := storage.RedisDB.Exists(context.TODO(), "detail_links:id:1").Val()
	log.Println("hasHome", hasHome)
	// 不存在首页的key 则认为是第一次启动
	if hasHome == 0 {
		spider.Tiankong().Start()
	}
}

func TimingSpider(cmd func()) {

	log.Println("cron TimingSpider start:")

	// v3 用法 干
	c := cron.New(cron.WithSeconds())

	// 每天定时执行的条件
	spec := viper.GetString(`cron.timing_spider`)

	c.AddFunc(spec, func() {
		//go StartSpider()
		//go spider.StartApi()
		cmd()
	})

	c.Start()

	// 关闭计划任务, 但是不能关闭已经在执行中的任务.
	// defer c.Stop()

	// 阻塞
	//select {}
}

func RecentUpdate(cmd func()) {
	c := cron.New()
	// 每天定时执行的条件
	spec := "0 12 * * *"

	c.AddFunc(spec, func() {
		cmd()
	})

	c.Start()
}
