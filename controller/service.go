package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type Service struct {
}

func ServiceRegister(group *gin.RouterGroup) {
	service := &Service{}
	//服务列表
	group.GET("/service_list", service.ServiceList)
	//删除服务
	group.GET("/service_delete", service.ServiceDelete)
	//获取单个服务
	group.GET("/service_detail", service.ServiceDetail)
	//服务流量统计
	group.GET("/service_stat", service.ServiceStat)
	//添加HTTP服务
	group.POST("/service_add_http", service.ServiceAddHTTP)
	//更新HTTP服务
	group.POST("/service_update_http", service.ServiceUpdateHTTP)

	/*group.POST("/service_add_tcp", service.serviceAddTcp)
	group.POST("/service_update_tcp", service.serviceUpdateTcp)
	group.POST("/service_add_grpc", service.serviceAddGrpc)
	group.POST("/service_update_grpc", service.serviceUpdateGrpc)*/

}

// ServiceList godoc
// @Summary 服务列表
// @Description 服务列表
// @Tags 服务管理
// @ID /service/service_list
// @Accept  json
// @Produce  json
// @Param info query string false "关键词"
// @Param page_size query int true "每页的条数"
// @Param page_num query int true "页数"
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /service/service_list [get]
func (s *Service) ServiceList(c *gin.Context) {
	params := &dto.ServiceListInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//业务逻辑开发
	serviceInfo := &dao.ServiceInfo{}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	list, total, err := serviceInfo.PageList(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
	}
	outList := []dto.ServiceListItemOutput{}
	for _, listItem := range list {
		serviceDetail, err := listItem.ServiceDetail(c, tx, &listItem)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			return
		}

		serviceAddr := "unkonw"
		clusterIp := lib.GetStringConf("base.cluster.cluster_ip")
		clusterPort := lib.GetStringConf("base.cluster.cluster_port")
		clusterSSLPort := lib.GetStringConf("base.cluster.cluster_ssl_port")

		//http后缀接入的方式： 需要使用集群的ip：clusterIP +clusterPort+path
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 1 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterSSLPort, serviceDetail.HTTPRule.Rule)

		}
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypePrefixURL &&
			serviceDetail.HTTPRule.NeedHttps == 0 {
			serviceAddr = fmt.Sprintf("%s:%s%s", clusterIp, clusterPort, serviceDetail.HTTPRule.Rule)

		}

		//http +域名的接入方式：domain
		if serviceDetail.Info.LoadType == public.LoadTypeHTTP && serviceDetail.HTTPRule.RuleType == public.HTTPRuleTypeDomain {
			serviceAddr = serviceDetail.HTTPRule.Rule

		}
		//tcp和grpc的接入方式:clusterIP+servicePort
		if serviceDetail.Info.LoadType == public.LoadTypeTCP {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.TCPRule.Port)
		}
		if serviceDetail.Info.LoadType == public.LoadTypeGRPC {
			serviceAddr = fmt.Sprintf("%s:%d", clusterIp, serviceDetail.GRPCRule.Port)
		}
		ipList := serviceDetail.LoadBalance.GetIPListByMode()
		outItem := dto.ServiceListItemOutput{
			ID:          listItem.Id,
			ServiceName: listItem.ServiceName,
			ServiceDesc: listItem.ServiceDesc,
			ServiceAddr: serviceAddr,
			Qps:         0,
			Qpd:         0,
			TotalNode:   len(ipList),
		}
		outList = append(outList, outItem)
	}
	out := &dto.ServiceListOutput{
		Total: total,
		List:  outList,
	}

	middleware.ResponseSuccess(c, out)
}

// ServiceDelete godoc
// @Summary 服务删除
// @Description 服务删除
// @Tags 服务管理
// @ID /service/service_delete
// @Accept  json
// @Produce  json
// @Param id query string true "服务的id"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /service/service_delete [get]
func (s *Service) ServiceDelete(c *gin.Context) {
	params := &dto.ServiceDeleteInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//标记删除
	serviceInfo.IsDelete = 1
	if err = serviceInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	middleware.ResponseSuccess(c, "")
}

// ServiceDetail godoc
// @Summary 获取单个服务
// @Description 获取单个服务
// @Tags 服务管理
// @ID /service/service_detail
// @Accept  json
// @Produce  json
// @Param id query string true "服务的id"
// @Success 200 {object} middleware.Response{data=dao.ServiceDetail} "success"
// @Router /service/service_detail [get]
func (s *Service) ServiceDetail(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}

	serviceInfo := &dao.ServiceInfo{Id: params.Id}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	serviceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, serviceDetail)
}

// ServiceStat godoc
// @Summary 服务流量统计
// @Description 服务流量统计
// @Tags 服务管理
// @ID /service/service_stat
// @Accept  json
// @Produce  json
// @Param id query string true "服务的id"
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /service/service_stat [get]
func (s *Service) ServiceStat(c *gin.Context) {
	params := &dto.ServiceDetailInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	/*	tx, err := lib.GetGormPool("default")
		if err != nil {
			middleware.ResponseError(c, 2001, err)
			return
		}*/

	/*	serviceInfo := &dao.ServiceInfo{Id: params.Id}
		serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
		if err != nil {
			middleware.ResponseError(c, 2002, err)
			return
		}*/
	todayList := []int64{}
	//当日流量统计
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yestertodayList := []int64{}
	for i := 0; i <= 23; i++ {
		todayList = append(todayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yestertodayList,
	})
}

//ServiceAddHTTP godoc
//@Summary 添加HTTP服务
//@Description 添加HTTP服务
//@Tags 服务管理
//@ID /service/service_add_http
//@Accept json
//@Produce json
//@Param body body dto.ServiceAddHTTPInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_add_http [post]
func (s *Service) ServiceAddHTTP(c *gin.Context) {
	params := &dto.ServiceAddHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//校验：ip和权重的数量是否相等
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("ip列表和权重列表的权重不一致"))
		return
	}
	//事务开启
	tx = tx.Begin()
	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	_, err = serviceInfo.Find(c, tx, serviceInfo)
	if err == nil {
		tx.Rollback()
		//建议使用："github.com/pkg/errors" 可以携带错误的堆栈
		middleware.ResponseError(c, 2003, errors.New("服务已存在"))
		return
	}
	//校验:接入类型
	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}

	//rule：1、接入类型为路径
	if _, err := httpUrl.Find(c, tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2004, errors.New("要接入的前缀已经存在"))
		return
	}

	//以事务的形式 将字段 加入到对应的库
	serviceModel := &dao.ServiceInfo{
		ServiceName: params.ServiceName,
		ServiceDesc: params.ServiceDesc,
	}
	if err := serviceModel.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}
	//规则表
	//主键 serviceModel.ID
	httpRule := &dao.HttpRule{
		ServiceID:      serviceModel.Id,
		RuleType:       params.RuleType,
		Rule:           params.Rule,
		NeedHttps:      params.NeedHttps,
		NeedStripUri:   params.NeedStripUri,
		NeedWebsocket:  params.NeedWebsocket,
		UrlRewrite:     params.UrlRewrite,
		HeaderTransfor: params.HeaderTransfor,
	}
	if err := httpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}
	//权限控制表
	accessControl := &dao.AccessControl{
		ServiceID:         serviceModel.Id,
		OpenAuth:          params.OpenAuth,
		BlackList:         params.BlackList,
		WhiteList:         params.WhiteList,
		ClientIPFlowLimit: params.ClientipFlowLimit,
		ServiceFlowLimit:  params.ServiceFlowLimit,
	}
	if err := accessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	//负载均衡表
	loadBalance := &dao.LoadBalance{
		ServiceID:              serviceModel.Id,
		RoundType:              params.RoundType,
		IpList:                 params.IpList,
		WeightList:             params.WeightList,
		UpstreamConnectTimeout: params.UpstreamConnectTimeout,
		UpstreamHeaderTimeout:  params.UpstreamHeaderTimeout,
		UpstreamIdleTimeout:    params.UpstreamIdleTimeout,
		UpstreamMaxIdle:        params.UpstreamMaxIdle,
	}
	if err := loadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}
	//提交事务
	tx.Commit()
	middleware.ResponseSuccess(c, "")
}

//ServiceUpdateHTTP godoc
//@Summary 修改HTTP服务
//@Description 修改HTTP服务
//@Tags 服务管理
//@ID /service/service_update_http
//@Accept json
//@Produce json
//@Param body body dto.ServiceUpdateHTTPInput true "body"
//@Success 200 {object} middleware.Response{data=string} "success"
//@Router /service/service_update_http [post]
func (s *Service) ServiceUpdateHTTP(c *gin.Context) {
	params := &dto.ServiceUpdateHTTPInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//校验:接入类型
	httpUrl := &dao.HttpRule{RuleType: params.RuleType, Rule: params.Rule}

	//rule：1、接入类型为路径
	if _, err := httpUrl.Find(c, tx, httpUrl); err == nil {
		tx.Rollback()
		middleware.ResponseError(c, 2002, errors.New("要接入的前缀已经存在"))
		return
	}
	//校验：ip和权重的数量是否相等
	if len(strings.Split(params.IpList, "\n")) != len(strings.Split(params.WeightList, "\n")) {
		tx.Rollback()
		middleware.ResponseError(c, 2003, errors.New("ip列表和权重列表的权重不一致"))
		return
	}

	//事务开启
	tx = tx.Begin()

	serviceInfo := &dao.ServiceInfo{ServiceName: params.ServiceName}
	serviceInfo, err = serviceInfo.Find(c, tx, serviceInfo)
	ServiceDetail, err := serviceInfo.ServiceDetail(c, tx, serviceInfo)
	if err != nil {
		tx.Rollback()
		//建议使用："github.com/pkg/errors" 可以携带错误的堆栈
		middleware.ResponseError(c, 2004, errors.New("服务已存在"))
		return
	}

	//服务信息表
	info := ServiceDetail.Info
	info.ServiceDesc = params.ServiceDesc
	if err := info.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2005, err)
		return
	}

	//规则表
	//主键 serviceModel.ID
	HttpRule := ServiceDetail.HTTPRule
	HttpRule.NeedHttps = params.NeedHttps
	HttpRule.NeedStripUri = params.NeedStripUri
	HttpRule.NeedWebsocket = params.NeedWebsocket
	HttpRule.UrlRewrite = params.UrlRewrite
	HttpRule.HeaderTransfor = params.HeaderTransfor
	if err := HttpRule.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2006, err)
		return
	}

	//权限控制表
	AccessControl := ServiceDetail.AccessControl
	AccessControl.OpenAuth = params.OpenAuth
	AccessControl.BlackList = params.BlackList
	AccessControl.WhiteList = params.WhiteList
	AccessControl.ClientIPFlowLimit = params.ClientipFlowLimit
	AccessControl.ServiceFlowLimit = params.ServiceFlowLimit
	if err := AccessControl.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2007, err)
		return
	}

	//负载均衡表
	LoadBalance := ServiceDetail.LoadBalance
	LoadBalance.RoundType = params.RoundType
	LoadBalance.IpList = params.IpList
	LoadBalance.WeightList = params.WeightList
	LoadBalance.UpstreamConnectTimeout = params.UpstreamConnectTimeout
	LoadBalance.UpstreamHeaderTimeout = params.UpstreamHeaderTimeout
	LoadBalance.UpstreamIdleTimeout = params.UpstreamIdleTimeout
	LoadBalance.UpstreamMaxIdle = params.UpstreamMaxIdle
	if err := LoadBalance.Save(c, tx); err != nil {
		tx.Rollback()
		middleware.ResponseError(c, 2008, err)
		return
	}

	//提交事务
	tx.Commit()

	middleware.ResponseSuccess(c, "")
}
