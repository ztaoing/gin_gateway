package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
)

type HttpRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名，url_prefix=url前缀"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domian表示域名，type=url_prefix=时表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"支持https 1支持"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"启用websocket 1=启用"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能，每行一个"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持add、delete、edit、"`
}

func (h *HttpRule) TableName() string {
	return "gateway_service_http_rule"
}

func (h *HttpRule) Find(c *gin.Context, tx *gorm.DB, search *HttpRule) (*HttpRule, error) {
	model := &HttpRule{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

func (h *HttpRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(public.GetGinTraceContext(c)).Save(h).Error; err != nil {
		return err
	}
	return nil
}

func (h *HttpRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceId int64) ([]HttpRule, int64, error) {
	var list []HttpRule
	var count int64
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(h.TableName()).Select("*")
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
