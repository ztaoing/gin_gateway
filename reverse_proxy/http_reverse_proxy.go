/**
* @Author:zhoutao
* @Date:2020/7/11 下午2:49
 */

package reverse_proxy

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy/load_balance"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func NewLoadBalanceReverseProxy(c *gin.Context, lb load_balance.LoadBalance, transport *http.Transport) *httputil.ReverseProxy {
	//请求协调者
	director := func(req *http.Request) {
		//todo Get（ 入参 ）
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil || nextAddr == "" {
			panic("get next addr failed")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic("parse addr failed")
		}

		//重新定义request
		//encoded query values, without '?'
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		//todo ?
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host

		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}

		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") {
			return nil
		}
		return nil
	}

	//错误回调
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		middleware.ResponseError(c, 999, err)
	}

	return &httputil.ReverseProxy{
		Director:       director,
		ModifyResponse: modifyFunc,
		ErrorHandler:   errFunc,
	}
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
