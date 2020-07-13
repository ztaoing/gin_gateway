/**
* @Author:zhoutao
* @Date:2020/7/13 上午10:09
 */

package tcp_proxy_router

import (
	"context"
	"fmt"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy"
	"github.com/go1234.cn/gin_scaffold/tcp_proxy_middleware"
	"github.com/go1234.cn/gin_scaffold/tcp_server"
	"log"
	"net"
)

var tcpServerList = []*tcp_server.TcpServer{}

//启动tcp方服务
func TcpServerRun() {
	serviceList := dao.ServiceManagerHandler.GetTcpServiceList()
	for _, v := range serviceList {
		temp := v
		log.Printf("[INFO-START] tcp_proxy_run %v\n", temp.TCPRule.Port)

		//启动1个tcp服务器
		go func(serviceDetails *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetails.TCPRule.Port)
			rb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetails)
			if err != nil {
				log.Fatalf("[INFO-ERR] GetLoadBalancer port %v err:%v\n", addr, err)
				return
			}

			//中间件
			router := tcp_proxy_middleware.NewTcpSliceRouter()
			router.Group("/").Use(
				tcp_proxy_middleware.TCPFlowCountMiddleware(),
				tcp_proxy_middleware.TCPFlowLimiterMiddleware(),
				tcp_proxy_middleware.TCPWhiteListMiddleware(),
				tcp_proxy_middleware.TCPBlackListMiddleware(),
			)

			coreFunc := func(c *tcp_proxy_middleware.TcpSliceRouterContext) tcp_server.TCPHandler {
				return reverse_proxy.NewTcpLoadBalanceReverseProxy(c, rb)
			}

			baseCtx := context.Background()
			baseCtx = context.WithValue(baseCtx, "service", serviceDetails)

			//构建路由及设置中间件
			routerHandler := tcp_proxy_middleware.NewTcpSliceRouterHandler(coreFunc, router)
			tcpServer := &tcp_server.TcpServer{
				Addr:    string(serviceDetails.TCPRule.Port),
				Handler: routerHandler,
				BaseCtx: baseCtx, //传递到下游的context
			}
			tcpServerList = append(tcpServerList, tcpServer)

			//执行
			//不是正常关闭的时候才会执行报错
			if err := tcpServer.ListenAndServe(); err != nil && err != tcp_server.ErrServerClosed {
				log.Printf("[INFO] ListenAndServe error:%v\n", err)
			}
		}(temp)
	}
}

//关闭tcp服务
func TcpServerStop() {
	for _, v := range tcpServerList {
		v.Close()
		log.Printf("[INFO-STOP]tcp_proxy_stop %v\n", v.Addr)
	}
}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("info tcpHandler \n"))
}
