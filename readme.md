# Aria

Aria名字来源于美剧《冰与火之歌》中史塔克家族的小女儿的名字Arya。活泼好动是个典型的假小子，不喜欢女红礼仪反而喜欢舞刀弄剑，被父亲评价为具有“狼血”。

![](img/aria.jpeg)

## 用途

用于自动生成微服务应用的框架代码, 基于  [go-kit](https://github.com/go-kit/kit) 设计思想。

## 编译

运行`cmd`目录下的`build.sh`，脚本会将微服务框架的源文件压缩打包为二进制文件，并写入`cmd/assets.go`文件中，之后通过`go build`将所有资源打包为一个单独的可执行文件`aria-<system>-<version>`.

## 使用

1.  获取`aria`可执行文件（通过上面的“编译”步骤自行编译，或向他人索取）
2.  将`aria`拷贝至某`$PATH`目录下
3.  执行`aria new <project_name>`创建微服务项目

命令执行成功预期输出如下：

```
$ aria new test_project

   _____             .___
  /  _  \   _______  |   | _____
 /  /_\  \  \_  __ \ |   | \__  \
/    |    \  |  | \/ |   |  / __ \_
\____|__  /  |__|    |___| (____  /
        \/                      \/

/.../gopath/src/test_project/endpoint/endpoints.go
/.../gopath/src/test_project/endpoint/middleware.go
/.../gopath/src/test_project/endpoint/production.go
/.../gopath/src/test_project/main.go
/.../gopath/src/test_project/models/doc.go
/.../gopath/src/test_project/protocol/http.go
/.../gopath/src/test_project/protocol/production/gen.sh
/.../gopath/src/test_project/protocol/production/production.pb.go
/.../gopath/src/test_project/protocol/production/production.proto
/.../gopath/src/test_project/service/interface.go
/.../gopath/src/test_project/service/middleware.go
/.../gopath/src/test_project/service/production.go
/.../gopath/src/test_project/transport/grpc.go
/.../gopath/src/test_project/transport/grpc_test.go
/.../gopath/src/test_project/transport/http.go

Successfully create new project [test_project] in your GOPATH( /.../gopath ).
```