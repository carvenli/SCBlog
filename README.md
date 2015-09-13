# Golang语言编写到博客程序

作者博客：<http://shuang.ca>

# 程序说明

本程序使用Golang开发，为作者自用博客程序

# 安装说明

* 安装Golang
* 安装Git
* 安装MongoDB

## 获取所使用的第三方库

```bash
go get github.com/astaxie/beego
go get github.com/axgle/mahonia
go get gopkg.in/mgo.v2
```

## 配置conf/app.conf

```
appname = SCBlog
httpport = 80
runmode = pro
sessionon = true

dbhost = 127.0.0.1
dbport = 27017
dbname = SCBlog
dbuser =
dbpass =
adminuser = shuangca
adminpass = shuang.ca
```

## 编译运行

进入SCBlog源码目录，运行如下命令

```bash
go build
./SCBlog
```

# 授权方式

本程序遵循MIT授权

# 程序反馈

<http://shuang.ca>