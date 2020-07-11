package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
)

type AdminController struct {
}

func AdminsRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	//更新密码
	group.POST("/change_pwd", adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Param body body dto.AdminLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router/admin/admin_info [get]
func (a *AdminController) AdminInfo(c *gin.Context) {

	//1、读取sessionKey 对应的key

	sess := sessions.Default(c)
	sessionInfo := sess.Get(public.AdminSessionInfoKey)
	//将session转换为string
	sessionInfoStr := sessionInfo.(string)
	//2、将取出的数据封装成结构体
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfoStr)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//输出参数
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "hi",
		Roles:        []string{},
	}
	//校验成功
	middleware.ResponseSuccess(c, out)
}

// ChangePwd godoc
// @Summary 修改密码
// @Description 修改密码
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router/admin/change_pwd [post]
func (a *AdminController) ChangePwd(c *gin.Context) {
	params := &dto.ChangePwdInput{}
	//校验参数
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//完成密码修改需要的步骤
	//读取session中的用户信息 sessionInfo
	//利用sessionInfo.id读取数据库 adminInfo

	sess := sessions.Default(c)
	sessionInfo := sess.Get(public.AdminSessionInfoKey)
	//将session转换为string
	sessionInfoStr := sessionInfo.(string)
	//2、将取出的数据封装成结构体
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessionInfoStr)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//从数据库中读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, &dao.Admin{UserName: adminSessionInfo.UserName})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//生成密码 salt+password 并保存到数据库
	saltPassword := public.SaltPassword(adminInfo.Salt, params.Password)
	adminInfo.Password = saltPassword
	if err = adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	//输出参数
	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "hi",
		Roles:        []string{},
	}
	//校验成功
	middleware.ResponseSuccess(c, out)
}
