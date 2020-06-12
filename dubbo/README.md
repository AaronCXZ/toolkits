###### 使用时的环境变量配置(包含两个必须配置)

1. CONF_CONSUMER_FILE_PATH : Consumer 端配置文件路径，使用 consumer 时必需。
2. CONF_PROVIDER_FILE_PATH：Provider 端配置文件路径，使用 provider 时必需。
3. APP_LOG_CONF_FILE ：Log 日志文件路径，必需。
4. CONF_ROUTER_FILE_PATH：File Router 规则配置文件路径，使用 File Router 时需要。

###### 注入服务
```go
// 客户端
func init() {
  config.SetConsumerService(userProvider)
}
// 服务端
func init() {
  config.SetProviderService(new(UserProvider))
}
```
 ###### 注入序列化描述
 ```go
hessian.RegisterJavaEnum(Gender(MAN))
hessian.RegisterJavaEnum(Gender(WOMAN))
hessian.RegisterPOJO(&User{})
```