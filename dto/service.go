package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
)

//服务列表
type ServiceListInput struct {
	Info     string `json:"info" form:"info" comment:"关键词" example:"" validate:""`              //关键词
	PageNum  int    `json:"page_num" form:"page_num" comment:"页数 example:"1" validate:""`       //页数
	PageSize int    `json:"page_size" form:"page_size" comment:"每页条数" example:"20" validate；""` //每页页数

}

//每条记录显示的信息
type ServiceListItemOutput struct {
	ID          int64  `json:"id" form:"id" comment:"id" example:"1" validate:""`                          //id
	LoadType    int    `json:"load_type" form:"load_type" comment:"负载类型" example:"1" validate:""`          //负载的类型
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名称" example:"服务" validate:""`   //服务的名称
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"描述" example:"描述" validate:""`     //描述
	ServiceAddr string `json:"service_addr" form:"service_addr" comment:"服务地址" example:"" validate:""`     //服务地址
	Qps         int64  `json:"qps" form:"qps" comment:"qps" example:"10" validate:""`                      //qps
	Qpd         int64  `json:"qpd" form:"qpd" comment:"qpd" example:"10" validate:""`                      //qpd
	TotalNode   int    `json:"total_node" form:"total_node" comment:"total_node" example:"10" validate:""` //节点数
}

//输出
type ServiceListOutput struct {
	Total int64 `json:"total" form:"total" comment:"总数" example:"56" validate:""` //总数
	List  []ServiceListItemOutput
}

//服务流量统计
type ServiceStatOutput struct {
	Today     []int64 `json:"today" form:"today" comment:"今日服务统计" example:"100" validate:""`
	Yesterday []int64 `json:"yesterday" form:"yesterday" comment:"昨日服务统计" example:"100" validate:""`
}

//服务删除
type ServiceDeleteInput struct {
	Id int64 `json:"Id" form:"id" comment:"服务id"  example:"56" validate:"required"`
}

//获取单个服务
type ServiceDetailInput struct {
	Id int64 `json:"Id" form:"id" comment:"服务id"  example:"56" validate:"required"`
}

//添加HTTP服务
type ServiceAddHTTPInput struct {
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                           //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required"`                                 //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                      //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`          //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`        //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`              //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换" example:"" validate:"valid_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                  //关键词
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                            //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                            //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流	" example:"" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"" validate:"min=0"`       //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                                //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"" validate:"required,valid_ipportlist"`                        //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weightlist"`               //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时, 单位s" example:"" validate:"min=0"`   //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时, 单位s" example:"" validate:"min=0"` //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间, 单位s" example:"" validate:"min=0"`       //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`                     //最大空闲链接数
}

//添加HTTP服务
type ServiceUpdateHTTPInput struct {
	ID          int64  `json:"id" form:"id" comment:"服务id" example:"" validate:"required,min=1"`                                 //服务id
	ServiceName string `json:"service_name" form:"service_name" comment:"服务名" example:"" validate:"required,valid_service_name"` //服务名
	ServiceDesc string `json:"service_desc" form:"service_desc" comment:"服务描述" example:"" validate:"required,max=255,min=1"`     //服务描述

	RuleType       int    `json:"rule_type" form:"rule_type" comment:"接入类型" example:"" validate:"max=1,min=0"`                           //接入类型
	Rule           string `json:"rule" form:"rule" comment:"接入路径：域名或者前缀" example:"" validate:"required"`                                 //域名或者前缀
	NeedHttps      int    `json:"need_https" form:"need_https" comment:"支持https" example:"" validate:"max=1,min=0"`                      //支持https
	NeedStripUri   int    `json:"need_strip_uri" form:"need_strip_uri" comment:"启用strip_uri" example:"" validate:"max=1,min=0"`          //启用strip_uri
	NeedWebsocket  int    `json:"need_websocket" form:"need_websocket" comment:"是否支持websocket" example:"" validate:"max=1,min=0"`        //是否支持websocket
	UrlRewrite     string `json:"url_rewrite" form:"url_rewrite" comment:"url重写功能" example:"" validate:"valid_url_rewrite"`              //url重写功能
	HeaderTransfor string `json:"header_transfor" form:"header_transfor" comment:"header转换" example:"" validate:"valid_header_transfor"` //header转换

	OpenAuth          int    `json:"open_auth" form:"open_auth" comment:"是否开启权限" example:"" validate:"max=1,min=0"`                  //关键词
	BlackList         string `json:"black_list" form:"black_list" comment:"黑名单ip" example:"" validate:""`                            //黑名单ip
	WhiteList         string `json:"white_list" form:"white_list" comment:"白名单ip" example:"" validate:""`                            //白名单ip
	ClientipFlowLimit int    `json:"clientip_flow_limit" form:"clientip_flow_limit" comment:"客户端ip限流	" example:"" validate:"min=0"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" form:"service_flow_limit" comment:"服务端限流" example:"" validate:"min=0"`       //服务端限流

	RoundType              int    `json:"round_type" form:"round_type" comment:"轮询方式" example:"" validate:"max=3,min=0"`                                //轮询方式
	IpList                 string `json:"ip_list" form:"ip_list" comment:"ip列表" example:"127.0.0.1:80" validate:"required,valid_ipportlist"`            //ip列表
	WeightList             string `json:"weight_list" form:"weight_list" comment:"权重列表" example:"" validate:"required,valid_weightlist"`               //权重列表
	UpstreamConnectTimeout int    `json:"upstream_connect_timeout" form:"upstream_connect_timeout" comment:"建立连接超时, 单位s" example:"" validate:"min=0"`   //建立连接超时, 单位s
	UpstreamHeaderTimeout  int    `json:"upstream_header_timeout" form:"upstream_header_timeout" comment:"获取header超时, 单位s" example:"" validate:"min=0"` //获取header超时, 单位s
	UpstreamIdleTimeout    int    `json:"upstream_idle_timeout" form:"upstream_idle_timeout" comment:"链接最大空闲时间, 单位s" example:"" validate:"min=0"`       //链接最大空闲时间, 单位s
	UpstreamMaxIdle        int    `json:"upstream_max_idle" form:"upstream_max_idle" comment:"最大空闲链接数" example:"" validate:"min=0"`                     //最大空闲链接数
}

//校验参数
func (s *ServiceListInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, s)
}

//校验参数
func (sl *ServiceDeleteInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, sl)
}

//校验参数
func (sl *ServiceDetailInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, sl)
}

//校验参数
func (sa *ServiceAddHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, sa)
}

//校验参数
func (sa *ServiceUpdateHTTPInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, sa)
}
