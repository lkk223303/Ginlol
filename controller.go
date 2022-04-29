package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// @Summary 初始
// @Description 帳號密碼
// @Tags Hello
// @Success 200 body main.IndexData
// @Router / [get]
func test(c *gin.Context) {

	data := new(IndexData)
	data.Title = "帳號"
	data.Content = "密碼"

}

// @Summary 說HALO
// @Tags Hello
// @Produce json
// @Param user path string true "名字"
// @Success 200 {string} string
// @Router /demo/v1/hello/{user} [get]
func hello(c *gin.Context) {
	name := c.Param("user")
	greet := "HALO " + name
	c.JSON(200, greet)
}

// @Summary 呈現登入頁面
// @Tags login
// @Produce text/html
// @Router /login [get]
func loginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// @Summary "帳號密碼輸入,如果沒有就新增"
// @Tags login
// @accept mpfd
// @Produce application/json.
// @Param user formData string true "User struct"
// @Param password formData string true "User struct"
// @Param password-again formData string true "User struct"
// @Success 200 {string} json "{"status": "You are logged in!"}"
// @Failure 401 {string} json "{"status": "unauthorized"}"
// @Failure 400 {string} json "{"error": err.Error()}"
// @Router /login [post]
func loginAuth(c *gin.Context) {
	var form User

	// 綁定User data 到form *User
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判斷使用者是否存在資料庫 是否帳號密碼正確 是否為admin, 如果沒有就新增,如果有給予登入和token
	if err, user := QueryUser(DB, form.User); err != nil {
		// 判斷使用者是否一二次密碼相同
		if form.Password == form.PasswordAgain {
			if err := InsertUser(DB, form.User, form.Password); err == nil {
				c.JSON(http.StatusOK, gin.H{
					"success": "註冊成功",
				})
			} else {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": err.Error(),
				})
			}
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "兩次密碼須相同",
			})
		}

	} else {
		// Login process...查詢到user須做密碼驗證
		// 實做產生jwt token
		// 將token 存在cookie
		// 使用者cookie存入 redis 5分鐘
		var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

		if user.Password == form.Password {
			redisCmd := RC.Set(form.User, token, 30*time.Second)
			if redisCmd.Err() != nil {
				fmt.Println("Set error: ", redisCmd.Err())
				return
			} else {
				fmt.Println("token SET !", token)
			}

			c.JSON(http.StatusOK, gin.H{
				"success": "Logged in!",
			})
		}
	}

	return
}
