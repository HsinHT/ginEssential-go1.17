package main

import (
	"example.com/ginessential/common"
	"github.com/gin-gonic/gin"
)

func main() {

	common.InitDB()

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")go
}
