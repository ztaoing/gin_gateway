/**
* @Author:zhoutao
* @Date:2020/7/12 上午7:07
 */

package public

import (
	"github.com/garyburd/redigo/redis"
	"github.com/go1234.cn/gin_scaffold/golang_common/lib"
)

func RedisConfPipline(pip ...func(c redis.Conn)) error {
	//创建redis conn
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		return err
	}
	defer c.Close()

	for _, f := range pip {
		f(c)
	}
	//将输出缓冲刷新到redis server
	c.Flush()
	return nil
}

//执行redis指令
func RedisConfDo(commandName string, args ...interface{}) (interface{}, error) {
	c, err := lib.RedisConnFactory("default")
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c.Do(commandName, args)
}
