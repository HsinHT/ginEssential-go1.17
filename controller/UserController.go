package controller

import (
	"log"
	"net/http"

	"example.com/ginessential/common"
	"example.com/ginessential/dto"
	"example.com/ginessential/model"
	"example.com/ginessential/response"
	"example.com/ginessential/util"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
		response.Response(ctx, http.StatusUnprocessableEntity, 421, nil, "手機號碼必須為11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密碼不能少於6位")
		return
	}
	// 如果名稱沒有傳，給一個10位的隨機字串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	// 判斷手機號碼是否存在
	if isTelephoneExist(DB, telephone) {
		response.Response(ctx, http.StatusUnprocessableEntity, 423, nil, "用戶已經存在")
		return
	}

	// 創建用戶
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密錯誤")
		return
	}

	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	// 返回結果
	response.Success(ctx, nil, "註冊成功")
}

func Login(ctx *gin.Context) {
	DB := common.GetDB()

	// 獲取參數
	telephone := ctx.PostForm("telephone")
	password := ctx.PostForm("password")

	// 數據驗證
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 421, nil, "手機號碼必須為11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密碼不能少於6位")
		return
	}

	// 判斷手機號是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 424, nil, "用戶不存在")
	}

	// 判斷密碼是否正確
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密碼錯誤")
		return
	}

	// 發放 token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系統異常")
		log.Printf("token generate error : %v", err)
		return
	}

	// 返回結果
	response.Success(ctx, gin.H{"token": token}, "登入成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "登入訊息")
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)

	return user.ID != 0
}
