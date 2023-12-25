package services

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"go-spider/package/storage"
	"strings"
)

const CategoriesKey = "categories"

var Json = jsoniter.ConfigCompatibleWithStandardLibrary

type Categories struct {
	Link string       `json:"link"`
	Name string       `json:"name"`
	Sub  []Categories `json:"sub"`
}

func AllCategoryData() []Categories {
	categories := AllCategory()

	var nav []Categories
	Json.Unmarshal([]byte(categories), &nav)

	return nav
}

// 获取url中的链接
func TransformId(Url string) string {
	UrlStrSplit := strings.Split(Url, "-id-")[1]
	return strings.TrimRight(UrlStrSplit, ".html")
}

func AllCategory() string {
	return storage.RedisDB.Get(context.TODO(), CategoriesKey).Val()
}
