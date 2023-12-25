package services

import (
	"context"
	"github.com/spf13/viper"
	"go-spider/package/storage"
	"strconv"
	"strings"
	"sync"
)

const paginateCacheKey = "paginate"

type MovieListStruct struct {
	Link      string `json:"link"`
	Name      string `json:"name"`
	Category  string `json:"category"`
	Cover     string `json:"cover"`
	UpdatedAt string `json:"updated_at"`
	Starring  string `json:"starring"`
	Quality   string `json:"quality"`
}

type SortMovieListStruct []MovieListStruct

var (
	mutex sync.Mutex
)

// 查询指定范围的数据
func MovieListsRange(key string, start, stop int64) []MovieListStruct {
	var data []MovieListStruct

	var movieKeyMap MovieListStruct

	sStart := strconv.FormatInt(start, 10)
	sStop := strconv.FormatInt(stop, 10)

	// field
	cacheKey := "movie_lists_key:" + key + ":start:" + sStart + ":stop:" + sStop
	cacheHashKey := paginateCacheKey

	movieList := storage.RedisDB.HGet(context.TODO(), cacheHashKey, cacheKey).Val()

	if movieList != "" {
		Json.Unmarshal([]byte(movieList), &data)
		return data
	}

	movieKeys, _ := storage.RedisDB.ZRevRange(context.TODO(), key, start, stop).Result()

	for _, val := range movieKeys {

		MovieDetail := MovieDetail(val)

		info := MovieDetail["info"].(map[string]string)
		detail := MovieDetail["detail"].(map[string]interface{})

		movieKeyMap.Name = info["name"]
		movieKeyMap.Link = info["link"]
		movieKeyMap.Cover = info["cover"]
		movieKeyMap.Quality = info["quality"]
		if detail["update"] != nil {
			movieKeyMap.UpdatedAt = detail["update"].(string)
		}
		if detail["type"] != nil {
			movieKeyMap.Category = detail["type"].(string)
		}
		if detail["starring"] != nil {
			movieKeyMap.Starring = detail["starring"].(string)
		}

		mutex.Lock()
		data = append(data, movieKeyMap)
		mutex.Unlock()

	}

	// 这个应该是在for外面调用才对，之前居然也不报错。。。。。。
	byteData, _ := Json.MarshalIndent(data, "", " ")
	storage.RedisDB.HSetNX(context.TODO(), cacheHashKey, cacheKey, string(byteData)).Err()

	return data
}

func TransformCategoryId(link string) string {
	return TransformId(link)
}

func MovieDetail(link string) map[string]interface{} {
	mutex.Lock() //对共享资源加锁
	defer mutex.Unlock()
	data := make(map[string]interface{})

	details := storage.RedisDB.Keys(context.TODO(), "movies_detail:"+link+"*").Val()
	//details := models.RangeSCanMoviesKey("movies_detail:" + link + "*")

	detail := make(map[string]string)
	if len(details) > 0 {
		detail = storage.RedisDB.HGetAll(context.TODO(), details[0]).Val()
	}

	if detail["name"] == "" {
		// 重新采集
		//detailId := TransformCategoryId(link)
		//spider.Create().PageDetail(detailId)
	}

	var kuYunMap []map[string]interface{}
	Json.Unmarshal([]byte(detail["kuyun"]), &kuYunMap)

	var ckm3u8Map []map[string]interface{}
	Json.Unmarshal([]byte(detail["ckm3u8"]), &ckm3u8Map)

	var downloadMap []map[string]interface{}
	Json.Unmarshal([]byte(detail["download"]), &downloadMap)

	var _detailMap map[string]interface{}
	Json.Unmarshal([]byte(detail["detail"]), &_detailMap)

	data["kuyun"] = kuYunMap
	data["ckm3u8"] = ckm3u8Map
	data["download"] = downloadMap
	data["detail"] = _detailMap

	delete(detail, "kuyun")
	delete(detail, "ckm3u8")
	delete(detail, "download")
	delete(detail, "detail")

	data["info"] = detail

	isFilm := "1"

	if len(_detailMap) > 0 {
		if strings.Index(_detailMap["type"].(string), "片") == -1 { // ...片  or  ... 剧
			isFilm = "0"
		}
	}

	data["is_film"] = isFilm

	return data
}

// 搜索影片
func SearchMovies(key string) []MovieListStruct {
	var (
		data     []MovieListStruct
		wg       sync.WaitGroup
		chanPool = make(chan int, 1000)
	)

	//movieKeys := models.FindMoviesKey("*" + ":movie_name:" + key + "*")
	movieKeys := RangeSCanMoviesKey("*" + ":movie_name:" + key + "*")

	for _, val := range movieKeys {
		wg.Add(1)
		go func(val string) {
			chanPool <- 1

			var movieKeyMap MovieListStruct
			MovieDetail := MovieDetail(TransformLink(val))

			info := MovieDetail["info"].(map[string]string)
			detail := MovieDetail["detail"].(map[string]interface{})

			movieKeyMap.Name = info["name"]
			movieKeyMap.Link = info["link"]
			movieKeyMap.Cover = info["cover"]
			movieKeyMap.Quality = info["quality"]
			if detail["update"] != nil {
				movieKeyMap.UpdatedAt = detail["update"].(string)
			}
			if detail["type"] != nil {
				movieKeyMap.Category = detail["type"].(string)
			}
			if detail["starring"] != nil {
				movieKeyMap.Starring = detail["starring"].(string)
			}

			mutex.Lock()
			data = append(data, movieKeyMap)
			mutex.Unlock()

			<-chanPool
			wg.Done()
		}(val)
	}
	wg.Wait()
	close(chanPool)
	return data
}

// 获取实际链接url
func TransformLink(Url string) string {
	UrlStrSplit := strings.Split(Url, "movies_detail:")[1]

	return strings.Split(UrlStrSplit, ":movie_name:")[0]
}

func MoviesRecommend() interface{} {
	recommend := viper.Get(`recommend`)
	c := new([]interface{})
	if recommend == nil {
		return *c
	}
	return recommend
}

func RangeSCanMoviesKey(key string) []string {
	var (
		all   []string
		i     uint64
		mutex sync.Mutex
		wg    sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		i = 0
		for {
			s, c, _ := storage.RedisDB.Scan(context.TODO(), i, key, 10000).Result()
			//log.Println("s c",s, c,i)
			// 游标为0，停止循环
			if c == 0 {
				for _, val := range s {
					mutex.Lock()
					all = append(all, val)
					mutex.Unlock()
				}
				break
			}

			for _, val := range s {
				mutex.Lock()
				all = append(all, val)
				mutex.Unlock()
			}
			i = c

		}
		wg.Done()
	}()
	wg.Wait()
	return all
}
