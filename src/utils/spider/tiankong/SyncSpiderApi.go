package tiankong

import (
	"context"
	"go-spider/package/storage"
	"go-spider/src/utils/spider"
	"log"
	"runtime"
	"strconv"
	"time"
)

type SyncSpiderApi struct {
	spider.SpiderTask
}

func (spiderApi *SyncSpiderApi) Start() {
	go StartSyncApi()
}

func (spiderApi *SyncSpiderApi) PageDetail(id string) {
	go Detail(id, 0)
}

func (spiderApi *SyncSpiderApi) DoRecentUpdate() {
	DoRecentUpdate()
}

func StartSyncApi() {
	allMoviesDoneKeyExists := storage.RedisDB.Exists(context.TODO(), "all_movies_done").Val()
	if allMoviesDoneKeyExists > 0 {
		return
	}
	synclist(1)
}

func synclist(pg int) {
	// 执行时间标记
	startTime := time.Now()
	//defer ants.Release()
	//antPool, _ := ants.NewPool(100)

	//_f := initFastHttp()

	catePageCounts := getCategoryPageCount()

	log.Println(catePageCounts)

	for _, catePageCount := range catePageCounts {
		//wg.Add(1)
		categoryId := catePageCount.categoryId
		PageCount := catePageCount.PageCount

		// 协程改成按顺序执行 不然请求过快会被限制
		actionList(categoryId, pg, PageCount)

		//antPool.Submit(func() {
		//	// 这里不能直接使用 catePageCount.categoryId 、catePageCount.PageCount
		//	// 在 submit 之前赋值变量传进来
		//	actionList(categoryId, pg, PageCount)
		//	wg.Done()
		//})

	}

	//wg.Wait()

	// 结束时间标记
	endTime := time.Since(startTime)

	ExecSecondsS := strconv.FormatFloat(endTime.Seconds(), 'f', 2, 64)
	ExecMinutesS := strconv.FormatFloat(endTime.Minutes(), 'f', 2, 64)
	ExecHoursS := strconv.FormatFloat(endTime.Hours(), 'f', 2, 64)

	log.Println("执行完成......")

	// 删除已缓存的页面
	go DelAllListCacheKey()

	// 全量 done -> set done 永久Redis 标识 -> new corntab every min ( done key exist && recent_update_key expire ) -> set recent_update_key 1h expire -> do recent 3h update
	// 一周进行一次全量爬取，资源网站的电影ID是会变的，fuck!!!
	storage.RedisDB.SetNX(context.TODO(), "all_movies_done", "done", time.Second*604800)
	log.Println("本次爬虫执行时间为：" + ExecSecondsS + "秒 \n 即" + ExecMinutesS + "分钟 \n 即" + ExecHoursS + "小时 \n " + runtime.GOOS)
	// 钉钉通知
}
