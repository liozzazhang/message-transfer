## Cloopy(小云)
[![Build status](https://prow.k8s.io/badge.svg?jobs=post-test-infra-bazel)](https://testgrid.k8s.io/sig-testing-misc#post-bazel)

Cloopy是一款用来集成企业微信发通知的软件，用Go语言开发，轻量高效。

## To start using Cloopy
#### Compile

```bash
# Mac下编译
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go

# Linux下编译
$ CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
$ CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go 

# Windows 下编译 Mac 和 Linux 64位可执行程序
SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build main.go

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build main.go
```
>GOOS：目标平台的操作系统（darwin、freebsd、linux、windows）
>GOARCH：目标平台的体系架构（386、amd64、arm）
>交叉编译不支持 CGO 所以要禁用它

#### Run
```bash
# 以在Linux运行为例
$ ./cloopy
# 默认监听在port:12345
# swagger: http://localhost:12345/apidocs/?url=http://localhost:12345/apidocs.json
```
#### Integration

1. Add Group Robot In WeChat Work Group
2. Get Robot Webhook Url
3. POST Request

url:
```bash
url="http://127.0.0.1:12345/cloopy/send"
```
body:
```json
{
  "msgtype": "text",
  "text": {
    "content": "test"
  }
}
```
request:
```bash
curl ${url} \
    -H 'Content-Type: application/json' \
    -d '
    {
      "msgtype": "text",
      "text": {
        "content": "Hello Cloopy"
      },
      "webhook_url": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx"
    }'
```

#### grafana webhook
默认已经支持grafana的数据结构.
只需要在grafana Alerting-Notification channels里添加类型为webhook的channel即可。

格式：

`http://10.66.17.96:12345/cloopy/grafana?webhook=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=xxx`
> Query参数目前只支持webhook,传入对应企业微信群机器人的webhook地址。

## Support
If you have questions, reach out to [liozza@163.com].