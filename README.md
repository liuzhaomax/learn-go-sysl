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

`make`

`go mod tidy`

`go mod vendor`

加入下面链接提供的代码
> https://sysl.io/docs/tutorial-codegen

`go run cmd/Petdemo/main.go config/config.yaml`

xargs.exe
> https://github.com/DoctorLai/BatchUtils/blob/master/xargs.cmd
