package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
	"github.com/go1234.cn/gin_scaffold/public"
	"io/ioutil"
	"time"
)

// 请求进入日志
func RequestInLog(c *gin.Context) {
	traceContext := lib.NewTrace()
	if traceId := c.Request.Header.Get("com-header-rid"); traceId != "" {
		traceContext.TraceId = traceId
	}
	if spanId := c.Request.Header.Get("com-header-spanid"); spanId != "" {
		traceContext.SpanId = spanId
	}

	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)

	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back

	lib.Log.TagInfo(traceContext, "_com_request_in", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	})
}

// 请求输出日志
func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")

	startExecTime, _ := st.(time.Time)
	public.ComLogNotice(c, "_com_request_out", map[string]interface{}{
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	})
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if lib.GetBoolConf("base.log.file_writer.on") {
			//开启的时候才记录日志文件
			RequestInLog(c)
			defer RequestOutLog(c)
		}

		c.Next()
	}
}
