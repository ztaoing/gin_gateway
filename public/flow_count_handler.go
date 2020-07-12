/**
* @Author:zhoutao
* @Date:2020/7/12 上午7:36
* 流量控制
 */

package public

import (
	"sync"
	"time"
)

var FlowCountHandler *FlowCounter

type FlowCounter struct {
	RedisFlowCountMap   map[string]*RedisFlowCountService
	RedisFlowCountSlice []*RedisFlowCountService
	Locker              sync.RWMutex
}

func NewFlowCounter() *FlowCounter {
	return &FlowCounter{
		RedisFlowCountMap:   map[string]*RedisFlowCountService{},
		RedisFlowCountSlice: []*RedisFlowCountService{},
		Locker:              sync.RWMutex{},
	}
}

func init() {
	FlowCountHandler = NewFlowCounter()
}

func (counter *FlowCounter) GetCounter(serviceName string) (*RedisFlowCountService, error) {
	for _, v := range counter.RedisFlowCountSlice {
		if v.AppID == serviceName {
			return v, nil
		}
	}
	//未发现 则新建
	newCounter := NewRedisFlowCountService(serviceName, time.Second*1)
	counter.RedisFlowCountSlice = append(counter.RedisFlowCountSlice, newCounter)

	counter.Locker.Lock()
	defer counter.Locker.Unlock()
	counter.RedisFlowCountMap[serviceName] = newCounter

	return newCounter, nil
}
