package dto

//大盘显示项
type PanelGroupDataOutput struct {
	ServiceNum      int64 `json:"service_num"`       //服务数量
	AppNum          int64 `json:"app_num"`           //应用数量
	CurrentQPS      int64 `json:"current_qps"`       //当时的QPS
	TodayRequestNum int64 `json:"today_request_num"` //今天的请求数

}

type DashServiceStatItemOutPut struct {
	Name     string `json:"name"`
	Value    int64  `json:"value"`
	LoadType int    `json:"load_type"`
}
type DashServiceStatOutput struct {
	Legend []string                    `json:"legend"`
	Data   []DashServiceStatItemOutPut `json:"data"`
}
