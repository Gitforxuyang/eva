# eva

基础框架

|组件|功能|备注|
|:-----| ----: | :----: |
|protoc-gen-eva|protobuf plugin，用以生成proto代码||
|evaCli|脚手架工具 用来生成代码||
|eva|基本库|基本库包含基本功能和若干插件，如 redis mongo等|


待做：
```
1.统一配置中心并配置即时生效
2.client的selector改成读取etcd 自己做服务注册、发现、负载
3.其它plugin
4.broker
5.api server自动注册并发现服务
6.单元测试
```
