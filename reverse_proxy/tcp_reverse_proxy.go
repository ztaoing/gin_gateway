/**
* @Author:zhoutao
* @Date:2020/7/13 下午5:02
 */

package reverse_proxy

import (
	"context"
	"github.com/go1234.cn/gin_scaffold/reverse_proxy/load_balance"
	"github.com/go1234.cn/gin_scaffold/tcp_proxy_middleware"
	"io"
	"log"
	"net"
	"time"
)

var defaultDialer = new(net.Dialer)

//tcp反向代理
type TcpReverseProxy struct {
	Addr                 string
	ctx                  context.Context //单次请求单独设置
	KeepAlivePeriod      time.Duration
	DialTimeout          time.Duration
	DialContextFunc      func(ctx context.Context, netWork, address string) (net.Conn, error)
	OnDialErrorFunc      func(src net.Conn, dstDialErr error)
	ProxyProtocolVersion int
}

func (tr *TcpReverseProxy) dialTimeout() time.Duration {
	if tr.DialTimeout > 0 {
		return tr.DialTimeout
	}
	//设置默认超时时间
	return time.Second * 10
}

func (tr *TcpReverseProxy) dialContext() func(ctx context.Context, network, address string) (net.Conn, error) {
	if tr.DialContextFunc != nil {
		return tr.DialContextFunc
	}
	//设置默认方法
	return (&net.Dialer{
		Timeout:   tr.DialTimeout,
		KeepAlive: tr.KeepAlivePeriod,
	}).DialContext
}

func (tr *TcpReverseProxy) keepAlivePeriod() time.Duration {
	if tr.KeepAlivePeriod != 0 {
		return tr.KeepAlivePeriod
	}
	return time.Minute
}

func (tr *TcpReverseProxy) onDialError() func(src net.Conn, destDialErr error) {
	if tr.OnDialErrorFunc != nil {
		return tr.OnDialErrorFunc
	}
	return func(src net.Conn, destDialErr error) {
		log.Printf("tcpproxy: for incoming conn %v,error dialing %q:%v", src.RemoteAddr().String(), tr.Addr, destDialErr)
		src.Close()
	}
}

//传入上游的conn，在这里完成与下游连接的数据交换
func (tr *TcpReverseProxy) ServeTCP(ctx context.Context, src net.Conn) {
	var cancel context.CancelFunc
	//设置连接超时
	if tr.DialTimeout >= 0 {
		ctx, cancel = context.WithTimeout(ctx, tr.DialTimeout)
	}

	dest, err := tr.DialContextFunc(ctx, "tcp", tr.Addr)

	if cancel != nil {
		cancel()
	}

	if err != nil {
		tr.OnDialErrorFunc(src, err)
		return
	}

	//最后关闭下游连接
	defer func() {
		go dest.Close()
	}()

	//在数据请求之前设置dest得的keepAlive参数
	if v := tr.keepAlivePeriod(); v > 0 {
		if c, ok := dest.(*net.TCPConn); ok {
			//开启长连接
			c.SetKeepAlive(true)
			c.SetKeepAlivePeriod(v)
		}
	}
	errChan := make(chan error, 1)

	//交换本地与目标服务的数据
	go tr.proxyCopy(errChan, dest, src)
	go tr.proxyCopy(errChan, src, dest)
}

func (tr *TcpReverseProxy) proxyCopy(errChan chan<- error, dest, src net.Conn) {
	_, err := io.Copy(dest, src)
	errChan <- err
}

func NewTcpLoadBalanceReverseProxy(c *tcp_proxy_middleware.TcpSliceRouterContext, lb load_balance.LoadBalance) *TcpReverseProxy {
	return func() *TcpReverseProxy {
		nextAddr, err := lb.Get("")
		if err != nil {
			log.Fatal("get next addr failed")
		}
		return &TcpReverseProxy{
			ctx:             c.Ctx,
			Addr:            nextAddr,
			KeepAlivePeriod: time.Second,
			DialTimeout:     time.Second,
		}
	}()
}
