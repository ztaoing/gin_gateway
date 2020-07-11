/**
* @Author:zhoutao
* @Date:2020/7/11 上午8:21
 */

package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

const (
	DefaultCheckMethod    = 0
	DefaultCheckMaxErrNum = 2
	DefaultCHeckTimeout   = 5
	DefaultCheckInterval  = 5
)

//负载均衡检测
type LoadBalanceCheckConf struct {
	observers    []ObServer
	confIpWeight map[string]string
	activeList   []string
	format       string
}

func (l *LoadBalanceCheckConf) Attach(o ObServer) {
	l.observers = append(l.observers, o)
}

//获取配置list
func (l *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range l.activeList {
		weight, ok := l.confIpWeight[ip]
		if !ok {
			//默认weight
			weight = "50"
		}
		confList = append(confList, fmt.Sprintf(l.format, ip)+","+weight)

	}
	return confList
}

//更新配置 通知监听者更新
func (l *LoadBalanceCheckConf) WatchConf() {
	go func() {
		confIpErrNum := map[string]int{}
		for {
			changedList := []string{}
			for item, _ := range l.confIpWeight {
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCHeckTimeout)*time.Second)
				if err == nil {
					//连接成功
					conn.Close()
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}

				if err != nil {
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] += 1
					} else {
						confIpErrNum[item] = 1
					}
				}
				//当前错误数 < 允许的最大错误数
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			sort.Strings(changedList)
			sort.Strings(l.activeList)
			//两者不相等
			if !reflect.DeepEqual(changedList, l.activeList) {
				l.UpdateConf(changedList)
			}
			//检测间隔
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

//更新配置 监听者也更新
func (l *LoadBalanceCheckConf) UpdateConf(conf []string) {
	fmt.Println("update config", conf)
	l.activeList = conf
	//监听者也更新
	for _, obs := range l.observers {
		obs.Update()
	}

}

func NewLoadBalanceConf(format string, confIpWeight map[string]string) (*LoadBalanceCheckConf, error) {
	aList := []string{}
	for item, _ := range confIpWeight {
		aList = append(aList, item)

	}
	CheckConf := &LoadBalanceCheckConf{format: format, activeList: aList, confIpWeight: confIpWeight}
	CheckConf.WatchConf()
	return CheckConf, nil
}
