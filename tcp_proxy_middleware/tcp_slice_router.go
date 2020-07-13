/**
* @Author:zhoutao
* @Date:2020/7/13 下午3:43
 */

package tcp_proxy_middleware

import (
	"context"
	"github.com/go1234.cn/gin_scaffold/tcp_server"
	"math"
	"net"
)

//最多63个中间件
const abortIndex int8 = math.MaxInt8

//router
type TcpSliceRouter struct {
	groups []*TcpSliceGroup
}

//创建group
func (g *TcpSliceRouter) Group(path string) *TcpSliceGroup {
	if path != "/" {
		panic("only accept path = /")
	}
	return &TcpSliceGroup{
		TcpSliceRouter: g,
		path:           path,
	}
}

func NewTcpSliceRouter() *TcpSliceRouter {
	return &TcpSliceRouter{}
}

//group
type TcpSliceGroup struct {
	*TcpSliceRouter
	path     string
	handlers []TcpHandlerFunc
}

type TcpHandlerFunc func(*TcpSliceRouterContext)

//构造回调方法
func (g *TcpSliceGroup) Use(middlewares ...TcpHandlerFunc) *TcpSliceGroup {
	g.handlers = append(g.handlers, middlewares...)

	exist := false
	//是否存在
	for _, oldGroup := range g.TcpSliceRouter.groups {
		if oldGroup == g {
			exist = true
		}
	}
	//不存在的时候添加到组中
	if !exist {
		g.TcpSliceRouter.groups = append(g.TcpSliceRouter.groups, g)
	}
}

/**
TcpSliceRouterContext
*/
type TcpSliceRouterContext struct {
	conn net.Conn
	Ctx  context.Context
	*TcpSliceGroup
	index int8
}

func (c *TcpSliceRouterContext) Get(key interface{}) interface{} {
	return c.Ctx.Value(key)
}

func (c *TcpSliceRouterContext) Set(key, val interface{}) {
	c.Ctx = context.WithValue(c.Ctx, key, val)
}

//最大值时跳出中间件
func (c *TcpSliceRouterContext) Abort() {
	c.index = abortIndex
}

func (c *TcpSliceRouterContext) IsAbort() bool {
	return c.index >= abortIndex
}

//重置回调
func (c *TcpSliceRouterContext) Reset() {
	c.index = -1
}

func (c *TcpSliceRouterContext) Next() {
	c.index++
	//后续还有未执行的handler func
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}
func newTcpSliceRouterContext(conn net.Conn, r *TcpSliceRouter, ctx context.Context) *TcpSliceRouterContext {
	newTcpsliceGroup := &TcpSliceGroup{}
	*newTcpsliceGroup = *r.groups[0]

	c := &TcpSliceRouterContext{
		conn:          conn,
		Ctx:           ctx,
		TcpSliceGroup: newTcpsliceGroup,
	}
	//重置回调
	c.Reset()
	return c
}

/**
TcpSliceRouterHandler
*/

type TcpSliceRouterHandler struct {
	coreFunc func(*TcpSliceRouterContext) tcp_server.TCPHandler
	router   *TcpSliceRouter
}

func NewTcpSliceRouterHandler(corefunc func(*TcpSliceRouterContext) tcp_server.TCPHandler, router *TcpSliceRouter) *TcpSliceRouterHandler {
	return &TcpSliceRouterHandler{
		coreFunc: corefunc,
		router:   router,
	}
}
func (t *TcpSliceRouterHandler) ServeTCP(ctx context.Context, conn net.Conn) {
	c := newTcpSliceRouterContext(conn, t.router, ctx)
	c.handlers = append(c.handlers, func(c *TcpSliceRouterContext) {
		t.coreFunc(c).ServeTCP(ctx, conn)
	})
	c.Reset()
	c.Next()
}
