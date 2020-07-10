package http_proxy_router

import (
	"context"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/cert_file"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"log"
	"net/http"
	"time"
)

var (
	HttpSrvHandler  *http.Server //80
	HttpsSrvHandler *http.Server //443
)

func HttpServerRun() {
	gin.SetMode(lib.GetStringConf("proxy.base.debug_mode"))
	//设置router 及中间件
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.http.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.http.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("proxy.http.max_header_bytes")),
	}

	log.Printf(" [INFO] http_proxy_run:%s\n", lib.GetStringConf("proxy.http.addr"))
	if err := HttpSrvHandler.ListenAndServe(); err != nil {
		log.Fatalf(" [ERROR] http_proxy_run:%s err:%v\n", lib.GetStringConf("proxy.http.addr"), err)
	}

}

func HttpServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] http_proxy_stopped\n")
}

//443
func HttpsServerRun() {
	gin.SetMode(lib.ConfBase.DebugMode)
	//设置router 及中间件
	r := InitRouter(middleware.RecoveryMiddleware(), middleware.RequestLog())
	HttpsSrvHandler = &http.Server{
		Addr:           lib.GetStringConf("proxy.https.addr"),
		Handler:        r,
		ReadTimeout:    time.Duration(lib.GetIntConf("proxy.https.read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(lib.GetIntConf("proxy.https.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << uint(lib.GetIntConf("base.https.max_header_bytes")),
	}

	log.Printf(" [INFO] HttpsServerRun:%s\n", lib.GetStringConf("proxy.https.addr"))
	//设置证书
	if err := HttpsSrvHandler.ListenAndServeTLS(cert_file.Path("server.crt"), cert_file.Path("server.key")); err != nil {
		log.Fatalf(" [ERROR] HttpsServerRun:%s err:%v\n", lib.GetStringConf("proxy.https.addr"), err)
	}

}

//443
func HttpsServerStop() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HttpsSrvHandler.Shutdown(ctx); err != nil {
		log.Fatalf(" [ERROR] HttpsServerStop err:%v\n", err)
	}
	log.Printf(" [INFO] HttpsServerStop stopped\n")
}
