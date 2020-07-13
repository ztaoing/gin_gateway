/**
* @Author:zhoutao
* @Date:2020/7/13 上午10:36
 */

package tcp_server

import (
	"context"
	"net"
)

/**
tcpKeepAliveListener
*/

type tcpKeepAliveListener struct {
	*net.TCPListener
}

//思考点：继承方法重写时，只要使用非指针接口
func (t tcpKeepAliveListener) Accept() (net.Conn, error) {
	tconn, err := t.AcceptTCP()
	if err != nil {
		return nil, err
	}
	return tconn, err

}

/**
contextKey
*/
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "tcp_proxy context value" + k.name
}

/**
conn
*/

var localAddrContextKey = &contextKey{"local-addr"}

type conn struct {
	server     *TcpServer
	cancelCtx  context.CancelFunc
	rwc        net.Conn
	remoteAddr string
}

func (c *conn) serve(ctx context.Context) {

	c.remoteAddr = c.rwc.RemoteAddr().String()
	ctx = context.WithValue(ctx, localAddrContextKey, c.rwc.LocalAddr())
	if c.server.Handler == nil {
		panic("handler empty")

	}
	c.server.Handler.ServeTCP(ctx, c.rwc)
}

func (c *conn) close() {
	c.rwc.Close()
}
