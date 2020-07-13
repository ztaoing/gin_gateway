/**
* @Author:zhoutao
* @Date:2020/7/12 下午1:41
 */

package dto

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/public"
)

type TokensInput struct {
	GrantType string `json:"grant_type" form"grant_type" comment:"授权类型" example:"client_credentials" validate:"required"` //授权类型
	Scope     string `json:"scope" form:"scope" comment:"权限范围 read write" example:"read_write" validate:"required"`       //权限范围

}

type TokensOutput struct {
	AccessToken string `json:"access_token" form:"access_token"` //访问令牌
	ExpiresIn   int    `json:"expires_in" form:"exipires_in"`    //过期时间
	TokenType   string `json:"token_type" form:"token_type"`     //token类型
	Scope       string `json:"scope" form:"scope"`               //权限范围
}

//绑定参数
func (param *TokensInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}
