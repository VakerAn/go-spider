package controller

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"go-spider/src/utils/spider/tiankong"
	"net/http"
)

func Debug(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	go tiankong.DelAllListCacheKey()

	fmt.Fprint(w, "DEBUG")
}
