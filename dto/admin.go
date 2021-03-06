package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
	"time"
)

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	LoginTime    time.Time `json:"login_time""`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

//修改密码
type ChangePwdInput struct {
	Password string `json:"password" form:"password" comment:"密码" example:"123456" validate:"required"`
}

//校验参数
func (param *ChangePwdInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
