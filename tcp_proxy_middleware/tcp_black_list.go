/**
* @Author:zhoutao
* @Date:2020/7/11 下午9:18
 */

package tcp_proxy_middleware

import (
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/public"
	"strings"
)

// ip 黑名单 访问控制
//白名单优先于黑名单，如果命中白名单就不会校验黑名单
func TCPBlackListMiddleware() func(c *TcpSliceRouterContext) {
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

		WhiteIpList := []string{}
		if serviceDetails.AccessControl.WhiteList != "" {
			WhiteIpList = strings.Split(serviceDetails.AccessControl.WhiteList, ",")
		}

		BlackIpList := strings.Split(serviceDetails.AccessControl.BlackList, ",")

		ClientString := c.conn.RemoteAddr().String() //ip:port

		splits := strings.Split(ClientString, ":")
		ClientIp := ""
		if len(splits) == 2 {
			ClientIp = splits[0]
		}
		//权限是否开启
		if serviceDetails.AccessControl.OpenAuth == 1 && len(WhiteIpList) == 0 && len(BlackIpList) > 0 {
			//当前的ip是否在白名单中
			if public.InStringSlice(BlackIpList, ClientIp) {
				//输出
				c.conn.Write([]byte("get none service"))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
