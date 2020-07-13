/**
* @Author:zhoutao
* @Date:2020/7/12 下午1:21
 */

package controller

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/pkg/errors"
	"strings"
	"time"
)

type OAuthController struct {
}

func OAuthRegister(g *gin.RouterGroup) {
	oauth := &OAuthController{}
	g.POST("/tokens", oauth.Tokens)
}

//token 处理
func (o *OAuthController) Tokens(c *gin.Context) {
	params := dto.TokensInput{}

	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}

	//获取认证
	splits := strings.Split(c.GetHeader("Authorization"), " ")
	if len(splits) != 2 {
		//1 basic 2用户名和密码
		middleware.ResponseError(c, 2001, errors.New("用户名或密码格式错误"))
		return
	}

	//解码:获取用户名和密码
	appIdSecret, err := base64.StdEncoding.DecodeString(splits[1])
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//取出 app_id+secret
	//生成 app_list缓存
	//匹配 app_id
	//基于jwt 生成token
	//生成output
	parts := strings.Split(string(appIdSecret), ":")
	if len(parts) != 2 {
		middleware.ResponseError(c, 2003, errors.New("用户名或密码格式错误"))
		return
	}

	//租户列表
	appList := dao.AppManagerHandler.GetApplist()
	for _, v := range appList {
		//验证用户名 验证密码
		if v.AppID == parts[0] && v.Secret == parts[1] {
			//生成声明
			claims := jwt.StandardClaims{
				Issuer:    v.AppID,
				ExpiresAt: time.Now().Add(public.JwtExpires * time.Second).In(lib.TimeLocation).Unix(),
			}
			//根据声明生成token
			token, err := public.JwtEncode(claims)
			if err != nil {
				middleware.ResponseError(c, 2004, err)
				return
			}
			output := &dto.TokensOutput{
				ExpiresIn:   public.JwtExpires,
				TokenType:   "bearer",
				AccessToken: token,
				Scope:       "read_write", //权限范围
			}
			middleware.ResponseSuccess(c, output)
			return

		}
	}
	middleware.ResponseError(c, 2005, errors.New("匹配APP信息失败"))

}
