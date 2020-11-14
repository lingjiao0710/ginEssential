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

https://github.com/lingjiao0710/ginEssential/commit/d9ea46bacc95a35348485a13b61451a8a19b3f29



## 2020-11-13 JWT配合中间件用户认证

安装jwt-go

```go
go get github.com/dgrijalva/jwt-go
```

在common包中新建jwt.go，增加ReleaseToken函数用于生成token，ParseToken用于解析token

```go
package common

import (
	"lingjiao0710/ginEssential/model"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_crect_string")

//Claims tocken结构体
type Claims struct {
	UserID uint
	jwt.StandardClaims
}

//ReleaseToken 生成token
func ReleaseToken(user model.User) (string, error) {
	//expirationTime token过期时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "lingjiao0710",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}
```

新建middleware包，增加AuthMiddleware.go用于添加中间件

```go
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		//校验token格式
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]

		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//验证通过后获取claim中的userID
		userID := claims.UserID
		db := common.GetDB()

		var user model.User
		db.First(&user, userID)

		//用户
		if user.ID == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "权限不足"})
			ctx.Abort()
			return
		}

		//用户存在 将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}
```

routes.go中增加中间件路由

```go
r.GET("api/auth/info", middleware.AuthMiddleware(), controller.Info)
```

UserController.go中新增Info，返回user数据

```go
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"code": 200,
			"data": gin.H{"user": user}})
}
```

修改Login.go支持发放token

```go
//判断密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "token系统异常"})
		log.Printf("token error :%v", err)
		return
	}

	//返回结果
	ctx.JSON(http.StatusOK, gin.H{
		"code":  200,
		"token": token,
		"msg":   "登录成功",
	})
```

修改完成后，先post一条login消息，复制生成的token，再获取info消息中填入复制的token，返回正确的user数据：

![image-20201113161847994](https://i.loli.net/2020/11/13/qBAnysjtZWUbm64.png)

![image-20201113161911247](https://i.loli.net/2020/11/13/62vYi7oy1H5ljDT.png)

本章代码

https://github.com/lingjiao0710/ginEssential/commit/31818002516688a1416056ee833e2e97606a8c4b