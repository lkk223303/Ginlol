package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)


type GoogleUser struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}


func CreateGoogleOAuth()*oauth2.Config{
	config:= &oauth2.Config{
		//憑證的 client_id
		ClientID:"1043531715683-rauouppcer4t41m8mh853r0qv31r676s.apps.googleusercontent.com",
		//憑證的 client_secret
		ClientSecret :"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",
		//當 Google auth server 驗證過後，接收從 Google auth server 傳來的資訊
		RedirectURL:  "http://localhost:8088/",
		//告知 Google auth server 授權範圍，在這邊是取得用戶基本資訊和Email，Scopes 為 Google 提供
   		Scopes: []string{
      		"https://www.googleapis.com/auth/userinfo.email",
	    	"https://www.googleapis.com/auth/userinfo.profile",
	    	},
	    //指的是 Google auth server 的 endpoint，用 lib 預設值即可
	    Endpoint: google.Endpoint,
	}
	
	// 這個 state 是為了防止 CSRF(跨站請求偽造) 攻擊而設置的。
	//可以隨機產出一個 state，當 Google server 驗證完後，會把 state 回傳給網站 server，如此一來我們就可以驗證 state以確保正確性。
	return config
}

// 這個會導向到google登入頁面
func LoginWithGoogleOAuth(c *gin.Context){
	config:= CreateGoogleOAuth()
	redirectURL:=config.AuthCodeURL("state")

	c.Redirect(http.StatusSeeOther,redirectURL) // http.StatusSeeOther 為 303
}


// 這個 state 是為了防止 CSRF(跨站請求偽造) 攻擊而設置
func checkStat(c *gin.Context){
	s:= c.Query("state")
	if s!= "state_parameter_passthrough_value"{
		c.AbortWithError(http.StatusUnauthorized,gin.Error{})
		return
	}
}

func GoogleLogin(c *gin.Context){
	code := c.Query("code")
	config := CreateGoogleOAuth()

	token, err := config.Exchange(oauth2.NoContext,code)
	if err != nil{
		log.Fatal(err)
		return
	}

	client := config.Client(oauth2.NoContext,token)
	res,err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err!=nil{
		c.AbortWithError(http.StatusUnauthorized,gin.Error{})
		return
	}
	defer res.Body.Close()
	
	// 讀取返回資料
	rawData,_:= ioutil.ReadAll(res.Body)

	// var user GoogleUser
	// json.Unmarshal(rawData,&user)
	// fmt.Println(user)
	name := gjson.GetBytes(rawData,"name").String()
	email := gjson.GetBytes(rawData,"email").String()

	fmt.Println("Email: ",email," user: ",name)
}

