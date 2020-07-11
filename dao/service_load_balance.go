package dao

import (
	"fmt"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy/load_balance"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type LoadBalance struct {
	ID            int64  `json:"id" gorm:"primary_key"`
	ServiceID     int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	CheckMethod   int    `json:"check_method" gorm:"column:check_method" description:"检查方法 tcpchk=检测端口是否握手成功	"`
	CheckTimeout  int    `json:"check_timeout" gorm:"column:check_timeout" description:"check超时时间	"`
	CheckInterval int    `json:"check_interval" gorm:"column:check_interval" description:"检查间隔, 单位s		"`
	RoundType     int    `json:"round_type" gorm:"column:round_type" description:"轮询方式 round/weight_round/random/ip_hash"`
	IpList        string `json:"ip_list" gorm:"column:ip_list" description:"ip列表"`
	WeightList    string `json:"weight_list" gorm:"column:weight_list" description:"权重列表"`
	ForbidList    string `json:"forbid_list" gorm:"column:forbid_list" description:"禁用ip列表"`

	UpstreamConnectTimeout int `json:"upstream_connect_timeout" gorm:"column:upstream_connect_timeout" description:"下游建立连接超时, 单位s"`
	UpstreamHeaderTimeout  int `json:"upstream_header_timeout" gorm:"column:upstream_header_timeout" description:"下游获取header超时, 单位s	"`
	UpstreamIdleTimeout    int `json:"upstream_idle_timeout" gorm:"column:upstream_idle_timeout" description:"下游链接最大空闲时间, 单位s	"`
	UpstreamMaxIdle        int `json:"upstream_max_idle" gorm:"column:upstream_max_idle" description:"下游最大空闲链接数"`
}

func (t *LoadBalance) TableName() string {
	return "gateway_service_load_balance"
}

func (t *LoadBalance) Find(c *gin.Context, tx *gorm.DB, search *LoadBalance) (*LoadBalance, error) {
	model := &LoadBalance{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (t *LoadBalance) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

func (t *LoadBalance) GetIPListByMode() []string {

	return strings.Split(t.IpList, ",")
}

func (t *LoadBalance) GetWeightListByMode() []string {

	return strings.Split(t.WeightList, ",")
}

var LoadBalancerHandler *LoadBalancer

//
type LoadBalancer struct {
	LoadBalanceMap   map[string]*LoadBalanceItem
	LoadBalanceSlice []*LoadBalanceItem
	Locker           sync.RWMutex
}
type LoadBalanceItem struct {
	LoadBalance load_balance.LoadBalance
	ServiceName string
}

func NewLoadBalancerHandler() *LoadBalancer {
	return &LoadBalancer{
		LoadBalanceMap:   map[string]*LoadBalanceItem{},
		LoadBalanceSlice: []*LoadBalanceItem{},
		Locker:           sync.RWMutex{},
	}
}

//获取loadbalance
func (l *LoadBalancer) GetLoadBalancer(service *ServiceDetail) (load_balance.LoadBalance, error) {
	//验证是否存在
	for _, lbItem := range l.LoadBalanceSlice {
		if lbItem.ServiceName == service.Info.ServiceName {
			return lbItem.LoadBalance, nil
		}
	}
	//前缀
	schema := "http"
	if service.HTTPRule.NeedHttps == 1 {
		schema = "https"
	}

	//前缀方式
	prefix := ""
	if service.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
		prefix = service.HTTPRule.Rule
	}

	ipList := service.LoadBalance.GetIPListByMode()
	weightList := service.LoadBalance.GetWeightListByMode()

	ipConf := map[string]string{}
	for k, v := range ipList {
		ipConf[v] = weightList[k]
	}

	Conf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s://%s", schema, prefix), ipConf)
	if err != nil {
		return nil, err
	}

	//生成负载均衡器
	lb := load_balance.LoadBalanceFactoryWithConf(load_balance.LbWeightRoundRobin, Conf)

	//保存到map和slice中
	lbItem := &LoadBalanceItem{
		LoadBalance: lb,
		ServiceName: service.Info.ServiceName,
	}

	l.LoadBalanceSlice = append(l.LoadBalanceSlice, lbItem)

	l.Locker.Lock()
	defer l.Locker.Unlock()
	l.LoadBalanceMap[service.Info.ServiceName] = lbItem
	return lb, nil

}

//连接池
//每一个服务使用一个单独的连接池

var TransportHandler *Transportor

type Transportor struct {
	TransportMap   map[string]*TransportItem
	TransportSlice []*TransportItem
	Locker         sync.RWMutex
}

//服务+对应的连接池
type TransportItem struct {
	Trans       *http.Transport
	ServiceName string
}

func NewTransportorHandler() *Transportor {
	return &Transportor{
		TransportMap:   map[string]*TransportItem{},
		TransportSlice: []*TransportItem{},
		Locker:         sync.RWMutex{},
	}
}

//根据服务获取对应的连接池
func (t *Transportor) GetTransportor(service *ServiceDetail) (*http.Transport, error) {
	//验证服务是否存在
	for _, transTtem := range t.TransportSlice {
		if transTtem.ServiceName == service.Info.ServiceName {
			//服务存在
			return transTtem.Trans, nil
		}
	}
	trans := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: time.Duration(service.LoadBalance.UpstreamConnectTimeout), //连接的超时时间

		}).DialContext,
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout),
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout),
	}
	//不存在时
	transItem := &TransportItem{
		Trans:       trans,
		ServiceName: service.Info.ServiceName,
	}
	t.TransportSlice = append(t.TransportSlice, transItem)
	t.Locker.Lock()
	defer t.Locker.Unlock()
	t.TransportMap[service.Info.ServiceName] = transItem

	return trans, nil
}

func init() {
	//负载均衡管理器
	//LoadBalancerHandler = NewLoadBalancerHandler()
	//连接池
	TransportHandler = NewTransportorHandler()
}
