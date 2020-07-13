/**
* @Author:zhoutao
* @Date:2020/7/12 下午2:32
* 租户管理
 */

package dao

import (
	"sync"
	"time"
)

/**
租户详情
*/
type App struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	AppID     string    `json:"app_id" gorm:"column:app_id" description:"租户id"`
	Name      string    `json:"name" gorm:"column:name" decsription:"租户名称"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"秘钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单 支持前缀匹配"`
	Qpd       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	Qps       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"created_at" gorm:"column:create_at" description:"添加时间"`
	UpdateAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否删除 0 否 1 是"`
}

/**
租户管理
*/
type AppManager struct {
	AppMap   map[string]*App
	AppSlice []*App
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

var AppManagerHandler *AppManager

func NewAppManager() *AppManager {
	return &AppManager{
		AppMap:   map[string]*App{},
		AppSlice: []*App{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

func init() {
	AppManagerHandler = NewAppManager()
}

func (a *AppManager) GetApplist() []*App {
	return a.AppSlice
}
