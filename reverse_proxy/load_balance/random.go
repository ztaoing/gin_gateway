/**
* @Author:zhoutao
* @Date:2020/7/11 上午9:21
 */

package load_balance

import (
	"errors"
	"math/rand"
)

type RandomBalance struct {
	curIndex int
	rss      []string
	//观察主体 在服务发现的负载均衡中用到
	conf LoadBalanceConf
}

//增加rss
func (r *RandomBalance) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("params lenght 1 at least")
	}

	for _, v := range params {
		addr := v
		r.rss = append(r.rss, addr)
	}

	return nil

}

func (r *RandomBalance) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	//随机取一个
	r.curIndex = rand.Intn(len(r.rss))
	return r.rss[r.curIndex]

}

func (r *RandomBalance) Get() (string, error) {
	return r.Next(), nil
}

func (r *RandomBalance) SetConf(conf LoadBalanceConf) {
	r.conf = conf
}

//更新
func (r *RandomBalance) Update() {
	panic("a ha!")
}
