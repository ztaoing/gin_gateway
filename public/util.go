package public

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

func SaltPassword(salt, password string) string {
	s1 := sha256.New()
	s1.Write([]byte(password))
	str1 := fmt.Sprintf("%x", s1.Sum(nil))

	//加salt
	s2 := sha256.New()
	s2.Write([]byte(str1 + salt))
	//格式化为16进制
	return fmt.Sprintf("%x", s2.Sum(nil))
}

//对象 到 json的 转换
func Obj2Json(s interface{}) string {
	data, _ := json.Marshal(s)
	return string(data)
}
