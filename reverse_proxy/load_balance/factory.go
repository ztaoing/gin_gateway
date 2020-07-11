/**
* @Author:zhoutao
* @Date:2020/7/11 上午8:58
 */

package load_balance

type LbType int

const (
	LbRandom LbType = iota
	LbRoundRobin
	LbWeightRoundRobin
	//LbConsistentHash
)

func LoadBanlanceFactory(types LbType) LoadBalance {
	switch types {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	default:
		return &RandomBalance{}
	}
}

func LoadBalanceWithConf(types LbType, conf LoadBalanceConf) LoadBalance {
	//观察者模式
	switch types {
	case LbRandom:
		lb := &RandomBalance{}
		lb.SetConf(conf)
		conf.Attach(lb)
		lb.Update()
		return lb

	case LbRoundRobin:
		lb := &RoundRobinBalance{}
		lb.SetConf(conf)
		conf.Attach(lb)
		lb.Update()
		return lb
	case LbWeightRoundRobin:
		lb := &WeightRoundRobinBalance{}
		lb.SetConf(conf)
		conf.Attach(lb)
		lb.Update()
		return lb
	default:
		lb := &RandomBalance{}
		lb.SetConf(conf)
		conf.Attach(lb)
		lb.Update()
		return lb
	}

}
