package dao

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
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
	Info          *ServiceInfo   `json:"info" description:"基本信息"`       //service基本信息
	HTTPRule      *HttpRule      `json:"http" description:""`           //http规则
	TCPRule       *TcpRule       `json:"tcp" description:""`            //tcp规则
	GRPCRule      *GrpcRule      `json:"grpc" description:""`           //grpc规则
	LoadBalance   *LoadBalance   `json:"load_balance" description:""`   //负载均衡规则
	AccessControl *AccessControl `json:"access_control" description:""` //访问控制规则 黑名单 白名单
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
			//取
			tmpItem := listItem
			serviceDetail, err := tmpItem.ServiceDetail(c, tx, &tmpItem)
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

	//选择ServiceSlice，无需加锁
	for _, serviceItem := range s.ServiceSlice {
		//是否是http服务
		if serviceItem.Info.LoadType != public.LoadTypeHTTP {
			//执行下一个
			continue
		}
		//如果是域名的匹配方式
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			//找到 匹配到域名
			if serviceItem.HTTPRule.Rule == host {
				return serviceItem, nil
			}
		}

		//域名前缀的设置
		if serviceItem.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL {
			//找到
			if strings.HasPrefix(path, serviceItem.HTTPRule.Rule) {
				return serviceItem, nil
			}
		}
	}

	return nil, errors.New("no matched service")
}

/*
tcp处理
*/

func (s *ServiceManager) GetTcpServiceList() []*ServiceDetail {
	list := []*ServiceDetail{}
	for _, v := range s.ServiceSlice {
		temp := v
		if temp.Info.LoadType == public.LoadTypeTCP {
			list = append(list, temp)
		}
	}
	return list
}

/**
grpc
*/

func (s *ServiceManager) GetGrpcServiceList() []*ServiceDetail {
	list := []*ServiceDetail{}
	for _, v := range s.ServiceSlice {
		temp := v
		if temp.Info.LoadType == public.LoadTypeTCP {
			list = append(list, temp)
		}
	}
	return list
}
