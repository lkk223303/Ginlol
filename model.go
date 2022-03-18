package main

/////////// Model define ///////////
type IndexData struct {
	Title   string
	Content string
}

// Login 會用到的 struct
type User struct {
	ID            string
	User          string `form:"user" binding:"required"`                   //使用者帳號
	Password      string `form:"password" binding:"required"`               //使用者密碼
	PasswordAgain string `form:"password-again" binding:"eqfield=Password"` //二次密碼
	Secret               //匿名欄位
}

type Admin struct {
	User
	AdminLevel int
}

type Secret struct {
	Token        string `json:"token"`
	GoogleOauth2 string `json:"googleoauth2"`
}
