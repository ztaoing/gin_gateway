package main

import (
	"flag"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/http_proxy_router"
	"github.com/go1234.cn/gin_scaffold/tcp_proxy_router"
	"os"
	"os/signal"
	"syscall"
)

var (
	//dashaboard后台管理 server代理服务器
	endpoint = flag.String("endpoint", "", "input endpoint dashaboard or server")
	//配置文件路径
	config = flag.String("config", "./conf/dev/", "input config file like ./conf/dev/")
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

		/*lib.InitModule("./conf/dev/")
		defer lib.Destroy()

		//start
		router.HttpServerRun()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		//stop
		router.HttpServerStop()*/
	} else {
		lib.InitModule(*config)

		defer lib.Destroy()

		//将服务加载到内存中
		dao.ServiceManagerHandler.LoadOnce()

		//start 80
		go func() {
			http_proxy_router.HttpServerRun()
		}()
		//start 443
		go func() {
			http_proxy_router.HttpsServerRun()
		}()

		//tcp服务
		go func() {
			tcp_proxy_router.TcpServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		//stop
		http_proxy_router.HttpServerStop()
		http_proxy_router.HttpsServerStop()
	}

}
