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

