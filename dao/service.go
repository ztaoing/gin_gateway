package dao

import (
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/public"
	"net/http/httptest"
	"strings"
	"sync"
)

var ServiceManagerHandler *ServiceManager

func init() {
	ServiceManagerHandler = NewServiceManager()
}

//服务详情
type ServiceDetail struct {
	Info          *ServiceInfo   `json:"info" description:"基本信息"` //基本信息
	HTTPRule      *HttpRule      `json:"http" description:""`
	TCPRule       *TcpRule       `json:"tcp" description:""`
	GRPCRule      *GrpcRule      `json:"grpc" description:""`
	LoadBalance   *LoadBalance   `json:"load_balance" description:""`
	AccessControl *AccessControl `json:"access_control" description:""`
}

//服务管理
type ServiceManager struct {
	ServiceMap   map[string]*ServiceDetail //服务名称：服务详情
	ServiceSlice []*ServiceDetail          //使用slice遍历获取
	Locker       sync.RWMutex
	init         sync.Once
	err          error
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		ServiceMap:   map[string]*ServiceDetail{},
		ServiceSlice: []*ServiceDetail{},
		Locker:       sync.RWMutex{},
		init:         sync.Once{},
	}
}

//将服务加载到内存中
func (s *ServiceManager) LoadOnce() error {
	//只执行一次
	s.init.Do(func() {
		//定义写入map和slice的方法

		//从db中分页读取基本信息
		seviceInfo := &ServiceInfo{}
		//模拟了一个responseWriter参数
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		//从连接池中获取
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.err = err
			return
		}
		params := &dto.ServiceListInput{PageNum: 1, PageSize: 999999}
		list, _, err := seviceInfo.PageList(c, tx, params)
		if err != nil {
			s.err = err
			return
		}

		s.Locker.Lock()
		defer s.Locker.Unlock()

		for _, listItem := range list {
			serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
			if err != nil {
				s.err = err
				return
			}
			//操作map的时候如果在设置的时候，读取的话，会造成无法找到内存溢出的情况
			s.ServiceMap[listItem.ServiceName] = serviceDetail
			s.ServiceSlice = append(s.ServiceSlice, serviceDetail)
		}

	})

	return s.err
}

//接入方式匹配
func (s *ServiceManager) HTTPAccessMode(c *gin.Context) (*ServiceDetail, error) {
	//匹配前缀 例如前缀是  /abc ==serviceSlice.rule 就匹配成功
	//域名匹配 www.text.com == serviceSlice.rule

	//host c.request.host
	//path c.request.url.path
	host := c.Request.Host //www.test.www:90900 同时包含了端口
	// :之前的域名
	host = host[0:strings.Index(host, ":")]

	path := c.Request.URL.Path

	matched := false
	//todo
	matchedService := &ServiceDetail{}
	//选择ServiceSlice，无需加锁
	for _, serviceItem := range s.ServiceSlice {
		//是否是http服务
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			//执行下一个
			continue
		}
		//如果是域名的匹配方式
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			if serviceItem.HTTPRule.Rule == host {
				//匹配到域名
				matched = true
				matchedService = serviceItem
			}
		}
	}

	return nil, nil
}
