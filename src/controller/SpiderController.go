package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-spider/src/utils/spider"
	"net/http"
)

func GoSpider(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	spider.Tiankong().Start()

	fmt.Fprint(w, "Spider ing....")
}
