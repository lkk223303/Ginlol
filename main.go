package main

import (
	"fmt"
	"net/http"

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
// @description Gin swagger
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8088
func main() {

	router := gin.Default()
	port := 8088
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port))
	router.LoadHTMLGlob("template/*")

	router.GET("/", test)
	router.Group("/demo/v1").GET("/hello/:user", hello)
	router.POST("/login", login)

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

// @Summary 初始
// @Description 登入
// @Tags Hello
// @accept mpfd
// @Produce json
// @Param form formData string true "帳號密碼輸入"
// @Success 200 {object} json
// @Router /login [post]
type Login struct {
	User     string `form:"user"`
	Password string `form:"password"`
}

func login(c *gin.Context) {
	var form Login
	fmt.Println(form.User, form.Password)
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if form.User == "admin" && form.Password == "123" {
		c.JSON(http.StatusOK, gin.H{
			"status": "You are logged in!",
		})
		return
	}
	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})

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
