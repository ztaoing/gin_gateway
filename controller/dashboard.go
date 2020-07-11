package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/dao"
	"github.com/go1234.cn/gin_scaffold/dto"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/middleware"
	"github.com/go1234.cn/gin_scaffold/public"
	"github.com/pkg/errors"
	"time"
)

type DashBoardController struct {
}

func DashBoardRegister(group *gin.RouterGroup) {
	service := &DashBoardController{}
	//指标统计
	group.GET("/panel_groupData", service.PanelGroupData)
	//流量统计
	group.GET("/flow_stat", service.FlowStat)
	//服务类型占比
	group.GET("/service_stat", service.ServiceStat)
}

// PanelGroupData godoc
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/panel_groupData
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceListOutput} "success"
// @Router /dashboard/panel_groupData [get]
func (s *DashBoardController) PanelGroupData(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	//服务数量
	_, serviceNum, err := serviceInfo.PageList(c, tx, &dto.ServiceListInput{PageSize: 1, PageNum: 1})
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	//APP 数量
	AppInfo := &dao.App{}
	_, APPNum, err := app.APPList(c, tx, &dto.APPListInput{PageNo: 1, PageSize: 1})
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	out := &dto.PanelGroupDataOutput{
		ServiceNum:      serviceNum,
		AppNum:          APPNum,
		CurrentQPS:      0,
		TodayRequestNum: 0,
	}

	middleware.ResponseSuccess(c, out)
}

// FlowStat godoc ppppkk
// @Summary 指标统计
// @Description 指标统计
// @Tags 首页大盘
// @ID /dashboard/flow_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.ServiceStatOutput} "success"
// @Router /dashboard/flow_stat [get]
func (s *DashBoardController) FlowStat(c *gin.Context) {

	todayList := []int64{}
	//当日流量统计
	for i := 0; i <= time.Now().Hour(); i++ {
		todayList = append(todayList, 0)
	}
	yestertodayList := []int64{}
	for i := 0; i <= 23; i++ {
		todayList = append(todayList, 0)
	}

	middleware.ResponseSuccess(c, &dto.ServiceStatOutput{
		Today:     todayList,
		Yesterday: yestertodayList,
	})
}

// ServiceStat godoc
// @Summary 服务统计
// @Description 服务统计
// @Tags 首页大盘
// @ID /dashboard/service_stat
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.DashServiceStatOutput} "success"
// @Router /dashboard/service_stat [get]
func (s *DashBoardController) ServiceStat(c *gin.Context) {
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	serviceInfo := &dao.ServiceInfo{}
	list, err := serviceInfo.GroupByLoadType(c, tx)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	legend := []string{}
	for index, item := range list {
		name, ok := public.LoadTypeMap[item.LoadType]
		if !ok {
			middleware.ResponseError(c, 2003, errors.New("LoadType is not found"))
			return
		}
		list[index].Name = name
		legend = append(legend, name)
	}
	out := &dto.DashServiceStatOutput{
		Legend: legend,
		Data:   list,
	}
	middleware.ResponseSuccess(c, out)
}
