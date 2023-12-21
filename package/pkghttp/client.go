package pkghttp

import (
	"bytes"
	"crypto/tls"
	"go-spider/package/zaplog"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var httpClient *http.Client

func InitHttp() {
	RoundTripper := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	httpClient = &http.Client{
		Transport: RoundTripper,
		Timeout:   time.Duration(15 * time.Second),
	}
}

func DoRequest(url string, method string, body string, headers map[string]string, retry ...int) (*http.Response, error) {
	start := time.Now()
	var err error
	var resp *http.Response
	var retryTimes = 1
	defer func() {
		cost := float64(time.Now().UnixMicro()-start.UnixMicro()) / 1000.0
		staus := -1
		respLen := -1
		if resp != nil {
			staus = resp.StatusCode
			respLen = int(resp.ContentLength)
		}
		zaplog.Logger.Info("HTTP SERVER",
			zap.Float64("COST", cost),
			zap.String("URL", url),
			zap.Any("header", headers),
			zap.Any("body", body),
			zap.Any("status_code", staus),
			zap.Any("resp_len", respLen),
			zap.Error(err))
	}()

	if len(retry) > 0 {
		retryTimes = retry[0]
	}
	for times := 1; times <= retryTimes; times++ {
		tempReq, e := http.NewRequest(method, url, bytes.NewBufferString(body))
		if e != nil {
			return nil, e
		}

		for key, value := range headers {
			tempReq.Header.Set(key, value)
		}
		resp, err = httpClient.Do(tempReq)
		if err == nil {
			// if there is no http request err , then break loop.
			break
		}
	}

	return resp, err
}
