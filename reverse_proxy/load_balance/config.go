/**
* @Author:zhoutao
* @Date:2020/7/11 上午8:21
 */

package load_balance

type LoadBalanceConf interface {
	Attach(o ObServer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

type ObServer interface {
	Update()
}
