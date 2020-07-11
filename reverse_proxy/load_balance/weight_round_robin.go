/**
* @Author:zhoutao
* @Date:2020/7/11 上午10:27
 */

package load_balance

import (
	"errors"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*Node
	rsw      []int
	conf     LoadBalanceConf
}

type Node struct {
	addr            string
	weight          int //权重值
	currentWeight   int //当前权重
	effectiveWeight int //有效权重
}

func (w *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		//第一个为目标服务器的ip，第二个为权重值
		return errors.New("params length need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}

	node := &Node{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight
	w.rss = append(w.rss, node)
	return nil
}

//获取加权后的地址
func (w *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *Node

	for i := 0; i < len(w.rss); i++ {
		r := w.rss[i]
		//统计所有有效权重之和
		total += r.effectiveWeight

		//变更节点的临时权重为节点的临时权重+节点有效权重
		r.currentWeight += r.effectiveWeight

		//有效权重默认与权重相同，通讯异常时为-1，通讯成功时为1，知道恢复到weight的大小
		if r.effectiveWeight < r.weight {
			r.effectiveWeight++
		}

		//选择最大临时权重节点
		if best == nil || r.currentWeight > best.currentWeight {
			best = r
		}

	}
	if best == nil {
		return ""
	}

	//变更临时权重= 临时权重 - 有效权重的和
	best.currentWeight -= total
	return best.addr
}

func (w *WeightRoundRobinBalance) Get(string) (string, error) {
	return w.Next(), nil
}

func (w *WeightRoundRobinBalance) SetConf(conf LoadBalanceConf) {
	w.conf = conf
}

func (w *WeightRoundRobinBalance) Update() {
	panic("a ha!")
}
