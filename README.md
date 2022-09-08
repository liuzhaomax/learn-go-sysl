# learn-go-sysl
learn-go-sysl

```shell
docker run -it --rm -v `pwd`:/work anzbank/sysl-go:latest specs/petdemo.sysl github.com/anz-bank/sysl-go-demo Petdemo:Petdemo
```

如果使用windows，上面命令中的`pwd`对于windows需要换成绝对路径
```shell
docker run -it --rm -v D:\workspace\GitHub\learn-go-sysl:/work anzbank/sysl-go:latest specs/petdemo.sysl github.com/anz-bank/sysl-go-demo Petdemo:Petdemo
```

如果使用windows，安装awk

如果使用windows，将`codegen.mk`的`$$(pwd)`换成`${CURDIR}`

如果使用windows，下载Cygwin。~~将xargs.exe加入环境变量。解决xargs无效问题。~~

> https://github.com/DoctorLai/BatchUtils/blob/master/xargs.cmd
> 此开源xargs.exe存在兼容性问题

如果使用windows，创建软连接，使linux可以找到windows文件
```shell
mklink /J "D:\workspace\GitHub\learn-go-sysl" "D:\cygwin64\home\LIUZHAO\sysl" 
```

如果使用windows，在cygwin的终端，软连接路径下，使用`make`命令

`make`

`go mod tidy`

`go mod vendor`

加入下面链接提供的代码
> https://sysl.io/docs/tutorial-codegen

`go run cmd/Petdemo/main.go config/config.yaml`

