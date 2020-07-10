package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/public"
	"time"
)

type ServiceInfo struct {
	Id          int64     `json:"id" gorm:"primary_key" description:"自增主键"`
	LoadType    int       `json:"load_type" gorm:"colum:load_type" description:"类型"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	IsDelete    int       `json:"is_delete" gorm:"columd:is_delete" description:"是否删除"`
}

func (d *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

//查找
func (d *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

//保存
func (d *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.SetCtx(public.GetGinTraceContext(c)).Save(d).Error

}

//
func (d *ServiceInfo) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutPut, error) {
	list := []dto.DashServiceStatItemOutPut{}

	//GetGinTraceContext保证链路中包含数据库的查询
	query := tx.SetCtx(public.GetGinTraceContext(c))
	if err := query.Table(d.TableName()).Where("is_delete=0").Select(
		"load_type ,count(*) as value").Group(
		"load_type").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (d *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	//计算偏移量
	offset := (param.PageNum - 1) * param.PageSize
	//GetGinTraceContext保证链路中包含数据库的查询
	query := tx.SetCtx(public.GetGinTraceContext(c))
	query = query.Table(d.TableName()).Where("is_delete=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%%")
	}
	//没有查询到记录
	if err := query.Limit(param.PageSize).Offset(offset).Order("id desc").Find(&list).Error; err != nil && err == gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	//总数
	query.Limit(param.PageSize).Offset(offset).Find(&list).Count(&total)
	return list, total, nil

}

//单个服务内容
func (d *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	//服务名称是否为空
	if search.ServiceName == "" {
		info, err := d.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	httpRule := &HttpRule{ServiceID: search.Id}
	httpRule, err := httpRule.Find(c, tx, httpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//tcp
	tcpRule := &TcpRule{ServiceID: search.Id}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//grpc
	grpcRule := &GrpcRule{ServiceID: search.Id}
	grpcRule, err = grpcRule.Find(c, tx, grpcRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//access
	access := &AccessControl{ServiceID: search.Id}
	access, err = access.Find(c, tx, access)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	//loadbalance
	load := &LoadBalance{ServiceID: search.Id}
	load, err = load.Find(c, tx, load)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      httpRule,
		TCPRule:       tcpRule,
		GRPCRule:      grpcRule,
		LoadBalance:   load,
		AccessControl: access,
	}
	return detail, nil
}
