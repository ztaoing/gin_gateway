package dao

import (
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
	"time"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"username" gorm:"colum:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"salt"`
	Password  string    `json:"password" gorm:"column:password" description:"密码"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"创建时间"`
	isDelete  int       `json:"is_delete" gorm:"columd:is_delete" description:"是否删除"`
}

func (d *Admin) TableName() string {
	return "gateway_admin"
}

//查找
func (d *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

//保存
func (d *Admin) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.SetCtx(public.GetGinTraceContext(c)).Save(d).Error

}
func (d *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := d.Find(c, tx, (&Admin{UserName: param.UserName, isDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	//构建saltpassword
	saltPassword := public.SaltPassword(adminInfo.Salt, param.Password)
	if adminInfo.Password != saltPassword {
		return nil, errors.New("密码错误，请重试")
	}
	return adminInfo, nil
}
