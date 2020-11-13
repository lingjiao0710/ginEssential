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



## 2020-11-13 实现用户登录

routes.go增加登录路由

```go
r.POST("api/auth/login", controller.Login)
```

UserContoller.go中实现Login函数

Register函数中保存密码需要进行加密

```go
//密码不能明文存储，需要对密码进行加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "密码加密错误"})
		return
	}
```

Login函数中判断密码并返回token

```go
//判断手机号是否存在
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}

	//判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token := 123

	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"token": token,
		"msg":   "登录成功",
	})
```

本章代码

