package main

import (
	"errors"
	"fmt"
	"net/http"

	_ "ginlol/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type IndexData struct {
	Title   string
	Content string
}

// @title Ginlol
// @version 1.0
// @description Ginlol
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8088
func main() {
	router := gin.Default() // 啟動server
	port := 8088            // port number setting
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port))
	router.LoadHTMLGlob("template/html/*")
	router.Static("/assets", "./template/assets")

	router.GET("/", test)
	router.Group("/demo/v1").GET("/hello/:user", hello)
	router.GET("/login", loginPage)
	router.POST("/login", loginAuth)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.Run(fmt.Sprintf(":%d", port))
}

// @Summary 初始
// @Description 帳號密碼
// @Tags Hello
// @Success 200 body main.IndexData
// @Router / [get]
func test(c *gin.Context) {

	data := new(IndexData)
	data.Title = "帳號"
	data.Content = "密碼"
	c.HTML(http.StatusOK, "index.html", data)
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

type Login struct {
	User          string `form:"user" binding:"required"`                   //判斷是否有輸入使用者帳號密碼
	Password      string `form:"password" binding:"required"`               //判斷是否有輸入使用者帳號密碼
	PasswordAgain string `form:"password-again" binding:"eqfield=Password"` //
}

// @Summary "帳號密碼輸入"
// @Tags login
// @accept mpfd
// @Produce application/json.
// @Param user formData string true "Login struct"
// @Param password formData string true "Login struct"
// @Param password-again formData string true "Login struct"
// @Success 200 {string} json "{"status": "You are logged in!"}"
// @Failure 401 {string} json "{"status": "unauthorized"}"
// @Failure 400 {string} json "{"error": err.Error()}"
// @Router /login [post]
func loginAuth(c *gin.Context) {
	var form Login

	// 綁定login data 到form *Login
	if err := c.Bind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 判斷使用者是否存在資料庫, 以及是否帳號密碼正確
	if err := Auth(form.User, form.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	} else if form.User == "admin" && form.Password == "123" { // 如果輸入 admin帳密
		c.JSON(http.StatusOK, gin.H{
			"success": "Admin logged in!",
		})
		return
	} else {
		// 判斷使用者是否一二次密碼相同
		if form.Password == form.PasswordAgain {
			c.JSON(http.StatusOK, gin.H{
				"success": "登入成功",
			})
			return

		} else { //不會進到這邊 因為在 struct Login 驗證
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errors.New("兩次輸入密碼不相同"),
			})
			return
		}
	}

}

// // @Summary 登入含大頭照 (050777)
// // @Description 帳密、大頭照上傳以及驗證
// // @Tags 代理紅利報表1
// // @Accept mpfd
// // @Produce json
// // @Param file formData file true "檔案"
// // @Param name path string true "帳號"
// // @Param password query string true "密碼"
// // @Param message formData string false "表單"
// // @Success 200 {object} json "{"fileName": file.Filename,"size": file.Size,"mimeType": file.Header,}"
// // @Failure 500 {object} json "{"error": err}"
// // @Router /agentBonus/RecordTestUpload/{name}/ [post]
// func RecordTestUpload(c *gin.Context) { // account, password, pic

// 	// 檔案上傳
// 	file, _, err := c.Request.FormFile("file")
// 	fmt.Printf("file: %b\n", file)
// 	if err != nil {
// 		// 	c.JSON(http.StatusInternalServerError, gin.H{
// 		// 		"error": err,
// 		// 	})
// 		fmt.Println("Error:", err)
// 	}

// 	// 參數跟Query
// 	id, password := c.Param("name"), c.Query("password")
// 	fmt.Println("Params:", id)
// 	fmt.Println("Query:", password)

// 	if password == "28" {
// 		c.String(200, "Welcome!, %s yr %s!", password, id)
// 		return
// 	}

// 	// 表單資料（Multipart/Urlencoded Form)
// 	message := c.PostForm("message")
// 	nickname := c.DefaultPostForm("nickname", "nickname")

// 	fmt.Println(message)
// 	c.JSON(200, gin.H{
// 		"message":  message,
// 		"nickname": nickname,
// 	})
// 	// 表單中單個欄位多個資料（Map as querystring or formPost parameters）

// 	// 上傳檔案（Upload Files）

// 	// 多個檔案（multiple files）
// 	// GET queryString / formData
// 	// Bind
// 	// var person Person
// 	// if err := c.Bind(&person); err != nil {
// 	// 	fmt.Println("Error: ", err)
// 	// }
// 	// fmt.Println(c.Bind(&person))

// 	// c.JSON(200, gin.H{
// 	// 	"name":     person.Name,
// 	// 	"password": person.Password,
// 	// })
// }

// type Person struct {
// 	Name     string `json:"name"`
// 	Password string `json:"password"`
// }
