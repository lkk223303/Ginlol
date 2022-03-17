package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "ginlol/docs"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// db const
var DB *sql.DB
var RC *redis.Client

const (
	USERNAME = "kent"
	PASSWORD = "0000"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "Ginlol"
)

/////////// Model define ///////////
type IndexData struct {
	Title   string
	Content string
}

// Login時會用到的 struct
type User struct {
	ID            string
	User          string `form:"user" binding:"required"`                   //使用者帳號
	Password      string `form:"password" binding:"required"`               //使用者密碼
	PasswordAgain string `form:"password-again" binding:"eqfield=Password"` //二次密碼
	Secret               //匿名欄位
}

// 新增user資料到資料庫(會檢查是否存在)
func (u *User) InsertUser(db *sql.DB, username, password string) error {

	fmt.Println("建立使用者中...")
	result, err := db.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("建立使用者失敗，原因：%v", err)
		return err
	} else {
		r, _ := result.RowsAffected()
		fmt.Println("建立資料成功 result:", r)

		return err
	}
}

type Admin struct {
	User
	AdminLevel int
}

type Secret struct {
	Token        string `json:"token"`
	GoogleOauth2 string `json:"googleoauth2"`
}
type Member interface {
	Register()
	Login()
	Delete()
	ChangePassword()
}

// type Changer interface {
// 	ChangePassword()
// }

// @title Ginlol
// @version 1.0
// @description Ginlol
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8088
func main() {

	// Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	defer client.Close()

	pong, err := client.Ping().Result()
	if err != nil {
		fmt.Println(err.Error())
		return
	} else {
		fmt.Println("ping result: ", pong)
	}
	RC = client

	// DB
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)

	dbConnect, err := sql.Open("mysql", conn)

	if err != nil {
		fmt.Println("開啟SQL連線發生錯誤: ", err)
		return
	}
	if err := dbConnect.Ping(); err != nil {
		fmt.Println("資料庫連線錯誤：", err)

	}

	DB = dbConnect
	DB.SetMaxOpenConns(10) // 可設置最大DB連線數，設<=0則無上限（連線分成 in-Use正在執行任務 及 idle執行完成後的閒置 兩種）
	DB.SetMaxIdleConns(10) // 設置最大idle閒置連線數。
	// 更多用法可以 進 sql.DBStats{}、sql.DB{} 裡面看
	defer dbConnect.Close()

	// Server
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

/////////// 資料庫操作 functions ///////////
// 對資料庫進行，建立一個名為 user 的 table
func CreateTable(db *sql.DB) error {
	sql := `CREATE TABLE IF NOT EXISTS users(id INT(4) PRIMARY KEY AUTO_INCREMENT NOT NULL,
	username VARCHAR(64),password VARCHAR(64)
	);`
	if _, err := db.Exec(sql); err != nil {
		fmt.Println("建立 Table發生錯誤: ", err)
		return err
	}
	fmt.Println("建立 Table 成功！")
	return nil
}

// 查詢user資料
func QueryUser(db *sql.DB, username string) (error, *User) {
	user := new(User)
	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	if err := row.Scan(&user.ID, &user.User, &user.Password); err != nil {
		fmt.Printf("查詢使用者失敗，原因：%v\n", err)
		return err, user
	}

	fmt.Printf("查詢使用者 ID: %s\tName: %s\tPassword: %s\t", user.ID, user.User, user.Password)
	return nil, user
}

// 新增user資料
func InsertUser(db *sql.DB, username, password string) error {

	fmt.Println("建立使用者中...")
	result, err := db.Exec("insert INTO users(username,password) values(?,?)", username, password)
	if err != nil {
		fmt.Printf("建立使用者失敗，原因：%v", err)
		return err
	} else {
		r, _ := result.RowsAffected()
		fmt.Println("建立資料成功 result:", r)

		return err
	}
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
	if err, _ := QueryUser(DB, form.User); err != nil {
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
		// 實做產生jwt tocken
		// 將token 存在cookie
		// 使用者cookie存入 redis 5分鐘
		var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

		// user := QueryUser(DB, form.User)
		// // 密碼驗證
		// if form.Password == user {
		// }

		err := RC.Set(form.User, token, 0)
		if err.Err() != nil {
			fmt.Println("Set error: ", err)
			return
		} else {
			fmt.Println("token SET !")
		}

		c.JSON(http.StatusOK, gin.H{
			"success": "Logged in!",
		})
		return
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
