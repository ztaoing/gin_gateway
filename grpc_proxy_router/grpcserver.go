/**
* @Author:zhoutao
* @Date:2020/7/13 上午10:09
 */

package grpc_proxy_router

import (
	"context"
	"fmt"
	"github.com/e421083458/grpc-proxy/proxy"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy"
	"google.golang.org/grpc"
	"log"
	"net"
)

var grpcServerList = []*warpGrpcServer{}

type warpGrpcServer struct {
	Addr string
	*grpc.Server
}

//启动tcp方服务
func GrpcServerRun() {
	serviceList := dao.ServiceManagerHandler.GetGrpcServiceList()
	for _, v := range serviceList {
		temp := v
		log.Printf("[INFO-START] tcp_proxy_run %v\n", temp.TCPRule.Port)

		//启动1个grpc服务器
		go func(serviceDetails *dao.ServiceDetail) {
			addr := fmt.Sprintf(":%d", serviceDetails.GRPCRule.Port)

			lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetails)
			if err != nil {
				log.Fatalf("[INFO-ERR] GetTcpLoadBalanceer %v err:%v\n", addr, err)
				return
			}

			lis, err := net.Listen("tcp", addr)
			if err != nil {
				log.Fatalf("[INFO] GrpcListen %v err:%v\n", addr, err)
				return
			}

			grpcHanler := reverse_proxy.NewGrpcLoadBalanceHandler(lb)

			s := grpc.NewServer(
				grpc.ChainStreamInterceptor(
					grpc_proxy_middleware.GrpcJwtAuthTokenMiddleware,
				),
				grpc.CustomCodec(proxy.Codec()),
				grpc.UnknownServiceHandler(grpcHanler),
			)

			grpcServerList = append(grpcServerList, &warpGrpcServer{
				Addr:   addr,
				Server: s,
			})

			log.Printf("[INFO] grpc_proxy_run %v\n", addr)

			if err := s.Serve(lis); err != nil {
				log.Fatalf("[INFO] grpc_proxy_run err:%v\n", err)
			}

		}(temp)
	}
}

//关闭tcp服务
func GrpcServerStop() {
	for _, grpcServer := range grpcServerList {
		grpcServer.GracefulStop()
		log.Printf("[INFO-STOP] grpc_proxy_stop %v\n", grpcServer.Addr, grpcServer.Addr)
	}
}

type tcpHandler struct {
}

func (t *tcpHandler) ServeTCP(ctx context.Context, src net.Conn) {
	src.Write([]byte("info tcpHandler \n"))
}
