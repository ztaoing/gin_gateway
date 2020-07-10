package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
)

//grpc规则
type GrpcRule struct {
	ID             int64  `json:"id" gorm:"primary_key"` //主键
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	Port           int    `json:"port" gorm:"column:port" description:"端口"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加 add、delete、edit"`
}

func (g *GrpcRule) TableName() string {
	return "gateway_service_grpc_rule"
}

func (g *GrpcRule) Find(c *gin.Context, tx *gorm.DB, search *GrpcRule) (*GrpcRule, error) {
	model := &GrpcRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (g *GrpcRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(g).Error; err != nil {
		return err
	}
	return nil
}

func (g *GrpcRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceId int64) ([]GrpcRule, int64, error) {
	var list []GrpcRule
	var count int64
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(g.TableName()).Select("*")
	query = query.Where("service_id=?", serviceId)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err

	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
