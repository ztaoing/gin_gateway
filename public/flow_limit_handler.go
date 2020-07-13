/**
* @Author:zhoutao
* @Date:2020/7/12 上午7:36
* 限流
 */

package public

import (
	"golang.org/x/time/rate"
	"sync"
)

var FlowLimitHandler *FlowLimiter

type FlowLimiter struct {
	FlowLimiterMap   map[string]*FlowLimiterItem
	FlowLimiterSlice []*FlowLimiterItem
	Locker           sync.RWMutex
}

//流量控制对象
type FlowLimiterItem struct {
	ServiceName string
	Limiter     *rate.Limiter
}

func NewFlowLimiter() *FlowLimiter {
	return &FlowLimiter{
		FlowLimiterMap:   map[string]*FlowLimiterItem{},
		FlowLimiterSlice: []*FlowLimiterItem{},
		Locker:           sync.RWMutex{},
	}
}

func init() {
	FlowLimitHandler = NewFlowLimiter()
}

func (counter *FlowLimiter) GetLimiter(serviceName string, qps float64) (*rate.Limiter, error) {
	for _, v := range counter.FlowLimiterSlice {
		if v.ServiceName == serviceName {
			return v.Limiter, nil
		}
	}
	//qps 每秒钟产生的token数
	newLimiter := rate.NewLimiter(rate.Limit(qps), int(qps*3))
	//未发现 则新建
	item := &FlowLimiterItem{
		ServiceName: serviceName,
		Limiter:     newLimiter,
	}

	counter.FlowLimiterSlice = append(counter.FlowLimiterSlice, item)

	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.FlowLimiterMap[serviceName] = item

	return newLimiter, nil
}
