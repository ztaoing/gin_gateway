/**
* @Author:zhoutao
* @Date:2020/7/11 上午7:58
* 反向代理中间件
 */

package http_proxy_middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy"
	"github.com/pkg/errors"
)

func HTTPReverseProxyMiddelware() gin.HandlerFunc {
	return func(c *gin.Context) {
		service, ok := c.Get("service")
		if !ok {
			middleware.ResponseError(c, 2001, errors.New("get services details failed"))
			//终止向下传递
			c.Abort()
			return
		}
		//服务
		serviceDetails := service.(*dao.ServiceDetail)
		//创建基于服务的负载均衡器
		lb, err := dao.LoadBalancerHandler.GetLoadBalancer(serviceDetails)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			c.Abort()
			return
		}
		//创建单个服务的连接池
		trans, err := dao.TransportHandler.GetTransportor(serviceDetails)
		if err != nil {
			middleware.ResponseError(c, 2003, err)
			c.Abort()
			return
		}
		//创建reverseProxy
		proxy := reverse_proxy.NewLoadBalanceReverseProxy(c, lb, trans)
		//使用reverseProxy.ServerHTTP(c.request,c,response)

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
