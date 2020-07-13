package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
)

type AccessControl struct {
	ID                int64  `json:"id" gorm:'primary_key'`
	ServiceID         int64  `json:"service_id" gorm:"column:service_id" description:"服务的id"`                      //服务的id
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip"`                      //黑名单id
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`                  //是否开启黑白名单权限 1=开启
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip"`                      //白名单ip
	WhiteHostName     string `json:"white_host_name" gorm:"column:white_host_name" description:"白名单主机"`            //主机白名单
	ClientIPFlowLimit int    `json:"client_ip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流"` //客户端ip限流
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流"`      //服务端限流
}

func (a *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

//查找
func (a *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	out := &AccessControl{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

//保存
func (a *AccessControl) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.SetCtx(public.GetGinTraceContext(c)).Save(a).Error

}
