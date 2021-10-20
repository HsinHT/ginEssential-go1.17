package controller

import (
	"log"
	"net/http"

	"example.com/ginessential/common"
	"example.com/ginessential/model"
	"example.com/ginessential/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	DB := common.GetDB()

	// 獲取參數
	name := ctx.PostForm("name")
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 數據驗證
	if len(telephone) != 11 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手機號碼必須為11位"})
		return
	}
	if len(password) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密碼不能少於6位"})
		return
	}
	// 如果名稱沒有傳，給一個10位的隨機字串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 判斷手機號碼是否存在
	if isTelephoneExist(DB, telephone) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用戶已經存在"})
		return
	}

	// 創建用戶
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	// 返回結果
	ctx.JSON(200, gin.H{
		"msg": "註冊成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)

	return user.ID != 0
}
