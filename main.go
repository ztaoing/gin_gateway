package main

import (
	"flag"
	"github.com/e421083458/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/http_proxy_router"
	"github.com/go1234.cn/gin_scaffold/router"
	"os"
	"os/signal"
	"syscall"
)

var (
	//dashaboard后台管理 server代理服务器
	endpoint = flag.String("endpoint", "", "input endpoint dashaboard or server")
	//配置文件路径
	config = flag.String("config", "", "input config file like ./conf/dev/")
)

func main() {
	flag.Parse()
	if *endpoint == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *config == "" {
		flag.Usage()
		os.Exit(1)
	}

	if *endpoint == "dashboard" {
		lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"})
		defer lib.Destroy()

		//start
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		//stop
		router.HttpServerStop()
	} else {
		lib.InitModule(*config, []string{"base", "mysql", "redis"})
		defer lib.Destroy()

		//start 80
		go func() {
			http_proxy_router.HttpServerRun()
		}()
		//start 443
		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		//stop
		http_proxy_router.HttpServerStop()
	}

}
