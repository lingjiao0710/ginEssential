package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

const (
	TELEPHONE_LEN    = 11
	PASSWORD_MIN_LEN = 6
)

//数据库表字段
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(110);not null;unique"`
	Password  string `gorm:"size:255;not null"`
}

func main() {

	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.POST("/api/auth/register", func(ctx *gin.Context) {

		//获取参数
		name := ctx.PostForm("name")
		telephone := ctx.PostForm("telephone")
		password := ctx.PostForm("password")

		//数据验证
		if len(telephone) != TELEPHONE_LEN {
			log.Println(len(telephone))
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位"})
			return
		}

		if len(password) < PASSWORD_MIN_LEN {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}

		//如果没有名称随机生成一个10位字符串
		if len(name) == 0 {
			name = RandomString(10)
		}
		log.Println(name, telephone, password)

		//判断手机号是否存在
		if isTelephoneExist(db, telephone) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户名已经存在"})
			return
		}

		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		//返回结果

		ctx.JSON(http.StatusOK, gin.H{
			"msg": "注册成功",
		})

	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

//isTelephoneExist 查询手机号是否存在
func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

//RandomString 生成一个随机字符串
func RandomString(n int) (s string) {
	var letters = []byte("asdf;lkjqwerpouiasdfjfjkalsdfjASDFQWERASDXCVB")
	result := make([]byte, n)

	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}

	return string(result)
}

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "root1234"
	charset := "utf8"

	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("连接数据库失败, error: " + err.Error())
	}

	//自动创建数据表
	db.AutoMigrate(&User{})

	return db

}
