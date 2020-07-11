/**
* @Author:zhoutao
* @Date:2020/7/11 上午9:36
 */

package load_balance

import "errors"

type RoundRobinBalance struct {
	curIndex int
	rss      []string
	conf     LoadBalanceConf
}

func (r *RoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("params len 1 at least")
	}

	r.rss = append(r.rss, params[0])
	return nil
}

func (r *RoundRobinBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}

	if r.curIndex >= len(r.rss) {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % len(r.rss)

	return curAddr
}

func (r *RoundRobinBalance) Get(string) (string, error) {
	return r.Next(), nil
}

func (r *RoundRobinBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

func (r *RoundRobinBalance) Update() {

}
