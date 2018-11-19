# nginx-monitor

收集nginx自带的http_stubs_tatus模块监控发送到falcon

## 系统环境

- running inside Linux
- nginx --with-http_stub_status_module

## 配置文件更改

1. 编辑conf.yaml文件，替换你的falcon地址和nginx状态检测地址
```
falcon-url: http://127.0.0.1:1988/v1/push
nginx-url: http://127.0.0.1:80/nginx_status
  
```
2. 需要第三方插件github.com/kylelemons/go-gypsy/yaml
```
go get github.com/kylelemons/go-gypsy/yaml

```
3. 启动方式
>\# crontab -e
  ```
* * * * * go run nginx_status.go

  ```

## Related Metrics

Metrics | Comments
--- | ---
Active | 对后端发起的活动连接数
accepts | Nginx总共处理的连接数
handled | 成功创建握手的数量
requests | 总共处理了多少请求
Reading | Nginx 读取到客户端的Header信息数
writing | Nginx 返回给客户端的Header信息数
waiting | 开启keep-alive后正在等候下一次请求指令的驻留连接
