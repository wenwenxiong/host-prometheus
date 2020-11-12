package apiserver

import (
	"github.com/wenwenxiong/host-prometheus/pkg/client/cache"
	"net/http"
)

type cacheHandler struct{
	client cache.Interface
}

func newCacheHandler( m cache.Interface) *cacheHandler {
	return &cacheHandler{ client: m}
}

func (chandler cacheHandler) MeticsCacheHandler(writer http.ResponseWriter, request *http.Request) {
	 hm := HostMetrics{}
	params := parseRequestParams(request)
	if params.host != "" {
		hm = getCache(params.host, chandler.client)
	}
	WriteAsJson(hm, writer)
}