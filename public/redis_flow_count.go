/**
* @Author:zhoutao
* @Date:2020/7/12 上午6:44
 */

package public

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"sync/atomic"
	"time"
)

//redis限流服务
type RedisFlowCountService struct {
	AppID       string
	Interver    time.Duration
	QPS         int64
	Unix        int64
	TickerCount int64
	TotalCount  int64
}

func NewRedisFlowCountService(appID string, interval time.Duration) *RedisFlowCountService {
	reqCounter := &RedisFlowCountService{
		AppID:    appID,
		Interver: interval,
		QPS:      0,
		Unix:     0,
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()

		ticker := time.NewTicker(interval)
		//定时同步
		for {
			<-ticker.C

			tickerCount := atomic.LoadInt64(&reqCounter.TickerCount)
			//重置
			atomic.StoreInt64(&reqCounter.TickerCount, 0)

			currentTime := time.Now()
			dayKey := reqCounter.GetDayKey(currentTime)
			hourKey := reqCounter.GetHourKey(currentTime)

			//以管道的形式发送
			if err := RedisConfPipline(func(c redis.Conn) {
				c.Send("INCRBY", dayKey, tickerCount)
				c.Send("EXPIRE", dayKey, 86400*2)

				c.Send("INCRBY", hourKey, tickerCount)
				c.Send("EXPIRE", hourKey, 86400*2)

			}); err != nil {
				fmt.Println("RedisConfPipline error:", err)
				continue
			}

			totalCount, err := reqCounter.GetDayData(currentTime)
			if err != nil {
				fmt.Println("reqCounter.GetDayData error:", err)
				continue
			}

			nowUnix := time.Now().Unix()
			if reqCounter.Unix == 0 {
				reqCounter.Unix = time.Now().Unix()
				continue
			}

			tickerCount = totalCount - reqCounter.TotalCount
			if nowUnix > reqCounter.Unix {
				reqCounter.TotalCount = totalCount
				reqCounter.QPS = tickerCount / (nowUnix - reqCounter.Unix)
				reqCounter.Unix = time.Now().Unix()
			}
		}
	}()
	return reqCounter
}

func (r *RedisFlowCountService) GetDayKey(t time.Time) string {
	dayStr := t.In(lib.TimeLocation).Format("20060102")
	return fmt.Sprintf("%s_%s_%s", RedisFlowDayKey, dayStr, r.AppID)
}

func (r *RedisFlowCountService) GetHourKey(t time.Time) string {
	hourStr := t.In(lib.TimeLocation).Format("2006010215")
	return fmt.Sprintf("%s_%s_%s", RedisFlowHourKey, hourStr, r.AppID)
}

//获取 天 key-value
func (r *RedisFlowCountService) GetDayData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", r.GetHourKey(t)))
}

func (r *RedisFlowCountService) GetHourData(t time.Time) (int64, error) {
	return redis.Int64(RedisConfDo("GET", r.GetHourKey(t)))
}

//原子增加
func (r *RedisFlowCountService) Increace() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		atomic.AddInt64(&r.TickerCount, 1)
	}()
}
