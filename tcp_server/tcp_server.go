/**
* @Author:zhoutao
* @Date:2020/7/13 上午10:36
 */

package tcp_server

import (
	"context"
	"fmt"
	"github.com/e421083458/go_gateway/tcp_server"
	"github.com/pkg/errors"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerClosed  = errors.New("tcp server closed")
	ServerContextKey = &contextKey{"tcp-server"}
)

type TcpServer struct {
	Addr    string
	Handler tcp_server.TCPHandler
	err     error
	BaseCtx context.Context

	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	KeepAliveTimeout time.Duration

	mu         sync.RWMutex
	inShutDown int32
	doneChan   chan struct{}
	onceClose  *onceCloseListener
}

//是否已经关闭
func (s *TcpServer) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutDown) != 0
}

func (s *TcpServer) ListenAndServe() error {
	//判断server是否已经关闭
	if s.shuttingDown() {
		return ErrServerClosed
	}

	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}

	addr := s.Addr
	if addr == "" {
		return errors.New("need addr")
	}

	//监听地址
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

func (s *TcpServer) Serve(l net.Listener) error {
	s.onceClose = &onceCloseListener{Listener: l}
	//关闭listener
	defer s.onceClose.Close()

	if s.BaseCtx == nil {
		s.BaseCtx = context.Background()
	}
	baseCtx := s.BaseCtx
	//将当前的TcpServer作为新的上下文内容
	ctx := context.WithValue(baseCtx, ServerContextKey, s)
	for {
		rw, err := l.Accept()
		if err != nil {
			select {
			case <-s.getDoneChan():
				return ErrServerClosed
			default:

			}
			fmt.Printf("accept failed ,error:%v\n", err)
			continue
		}
		c := s.newConn(rw)
		go c.serve(ctx)
	}
}

func (s *TcpServer) Close() error {
	atomic.StoreInt32(&s.inShutDown, 1)
	close(s.doneChan)
	s.Close() //关闭listener
	return nil
}

func (s *TcpServer) getDoneChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.doneChan == nil {
		s.doneChan = make(chan struct{})
	}
	return s.doneChan
}

func (s *TcpServer) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: s,
		rwc:    rwc,
	}

	//设置参数
	if d := c.server.ReadTimeout; d != 0 {
		c.rwc.SetReadDeadline(time.Now().Add(d))
	}
	if d := c.server.WriteTimeout; d != 0 {
		c.rwc.SetWriteDeadline(time.Now().Add(d))
	}
	if d := c.server.KeepAliveTimeout; d != 0 {
		if tcpConn, ok := c.rwc.(*net.TCPConn); ok {
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(d)
		}
	}

	return c
}

/**
onceCloseListener
*/
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

//关闭tcp listener
func (o *onceCloseListener) Close() error {
	//只关闭一次
	o.once.Do(o.close)
	return o.closeErr
}
func (o *onceCloseListener) close() {
	o.closeErr = o.Listener.Close()
}

/**
ListenAndServe
*/

type TCPHandler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}

func ListenAndServe(addr string, handler TCPHandler) error {
	server := &TcpServer{Addr: addr, Handler: handler, doneChan: make(chan struct{})}
	return server.ListenAndServe()
}
