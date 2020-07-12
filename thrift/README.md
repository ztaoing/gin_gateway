#使用thrift构建server和client
1. 编写 `thrift_gen.thrift`
2. 使用IDL命令生成server和client
`thrift --gen go thrift_gen.thrift`