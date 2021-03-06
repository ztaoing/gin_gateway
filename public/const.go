package public

const (
	ValidatorKey        = "ValidatorKey"
	TranslatorKey       = "TranslatorKey"
	AdminSessionInfoKey = "AdminSessionInfoKey"

	//负载类型
	LoadTypeHTTP = 0
	LoadTypeTCP  = 1
	LoadTypeGRPC = 2

	//接入方式
	HTTPRuleTypePrefixURL = 0
	HTTPRuleTypeDomain    = 1

	RedisFlowDayKey  = "fload_day_count"
	RedisFlowHourKey = "flow_hour_count"

	//统计
	FloatTotal        = "float_total"
	FloatCountService = "flow_service"
	FloatCountAPP     = "flow_app"
)

var (
	LoadTypeMap = map[int]string{
		LoadTypeHTTP: "HTTP",
		LoadTypeTCP:  "TCP",
		LoadTypeGRPC: "GRPC",
	}
)
