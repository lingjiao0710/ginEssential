package main

import (
	"lingjiao0710/ginEssential/common"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db := common.InitDB()
	defer db.Close()

	r := gin.Default()
	r = CollectRoute(r)
	r.Run() // listen and serve on 0.0.0.0:8080
}
