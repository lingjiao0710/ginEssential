# ginEssential
Gin + Vue 前后端分离项目实战



## 2020-11-12 实现用户注册

在github上新建项目，通过vscode克隆到本地

![image-20201112142549358](https://i.loli.net/2020/11/12/ky2ST3XBwNR4E5G.png)

新建main.go，并生成go.mod，同时下载gin包依赖

```shell
go mod init lingjiao0710/ginEssential
go get -u github.com/gin-gonic/gin
```

![image-20201112143850611](https://i.loli.net/2020/11/12/wNxrl1dO8EJPTIa.png)

到gin的官网https://gin-gonic.com/docs/quickstart/，复制示例程序到代码中验证gin服务：

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
```

修改api服务为用户注册，并添加用户注册代码，使用postman测试：

![image-20201112152327107](https://i.loli.net/2020/11/12/ryVitIvqOoHJ6RQ.png)

使用gorm（https://github.com/jinzhu/gorm）连接数据库，引入mysql数据库

```shell
go get -u github.com/jinzhu/gorm
go get -u github.com/go-sql-driver/mysql
```

本章代码

https://github.com/lingjiao0710/ginEssential/commit/7e36286ce380a2570d70441f8f07b8d3b658c047



## 2020-11-12 项目重构关注分离

将main.go的代码分离到各个包中，使代码模块化：

![image-20201112213322291](https://i.loli.net/2020/11/12/AMueOx4GHlVdyFX.png)

部分说明：

路由相关代码均移到routes.go中

将RandomString移到util包utils.go中

数据表定义移到model包user.go中

Register和isTelephoneExist函数移到controller包UserContoller.go中

数据库操作移到common包database.go中

修改代码时需要注意数据库的初始化和DB文件的获取，部分变量需要增加对应包名

编译和运行命令：

```shell
go build
ginEssential.exe
```

本章代码

https://github.com/lingjiao0710/ginEssential/commit/96d0cf5d48b46555e8b69fd89a05649cb30e0015