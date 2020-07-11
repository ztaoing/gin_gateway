/**
* @Author:zhoutao
* @Date:2020/7/11 上午9:23
 */

package load_balance

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	//后期服务更新
	Update()
}
