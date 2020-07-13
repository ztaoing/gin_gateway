/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package tcp_proxy_middleware

import (
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/public"
)

// 限流控制
func TCPFlowCountMiddleware() func(c *TcpSliceRouterContext) {
	return func(c *TcpSliceRouterContext) {
		service := c.Get("service")
		if service == nil {
			//没有拿到service的设置
			//不允许继续向下传递，即不再执行后边的中间件
			//输出
			c.conn.Write([]byte("get none service"))
			c.Abort()

			return
		}
		serviceDetails := service.(*dao.ServiceDetail)

		//全站的统计

		totalCount, err := public.FlowCountHandler.GetCounter(public.FloatTotal)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return

		}
		totalCount.Increace()

		//服务的统计
		ServiceCount, err := public.FlowCountHandler.GetCounter(public.FloatCountService + serviceDetails.Info.ServiceName)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		ServiceCount.Increace()

		//租户的统计
		APPCount, err := public.FlowCountHandler.GetCounter(public.FloatCountAPP)
		if err != nil {
			c.conn.Write([]byte(err.Error()))
			c.Abort()
			return
		}
		APPCount.Increace()

		c.Next()
	}
}
