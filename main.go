package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"net/http"
	_ "net/http/pprof"

	_ "Ginlol/docs"

	_ "github.com/gin-contrib/cors"
	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
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

// @title Ginlol
// @version 1.0
// @description Ginlol
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8088
func main() {

	// MariaDBOn()
	// RedisOn()

	// // gRPC init
	// // client := grpcClientOn()
	// // PB := pb.NewHelloServiceClient(client) //建立service通道
	// // defer grpcClientOff(client)

	// // RouteChat(PB,20)

	// // Profiling
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// // Server
	router := gin.Default() // 啟動server
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))

	port := 8888 // port number setting
	url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port))
	router.LoadHTMLGlob("template/html/*")
	router.Static("/assets", "./template/assets")

	router.GET("/", GoogleLogin)
	router.Group("/demo/v1").POST("/hello/:user/:gender/:gender", hello)
	router.GET("/login", loginPage)
	router.POST("/login", loginAuth)
	router.Group("/oauth").GET("/google/url", LoginWithGoogleOAuth)
	router.Group("/oauth").GET("/google/login", GoogleLogin)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.Run(fmt.Sprintf(":%d", port))

}
func swticher(i interface{}) {
	switch i.(type) {
	case int:
		fmt.Println("int: ", fmt.Sprintf("%d", i))

	case float32:
		fmt.Println("f32: ", fmt.Sprintf("%.2f", i))

	case float64:
		fmt.Println("f64: ", fmt.Sprintf("%.2f", i))

	case string:
		fmt.Println("str: ", fmt.Sprintf("%s", i))

	}
}

func channelControlServer() {
	// 用于监听服务退出
	done := make(chan error, 2)
	// 用于控制服务退出，传入同一个 stop，做到只要有一个服务退出了那么另外一个服务也会随之退出
	stop := make(chan struct{}, 0)
	// for debug
	go func() {
		pprof(stop)
	}()
	go func() {
		for i := 0; i < 5; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println(i)
		}
		done <- fmt.Errorf("It's an error")
	}()
	// main service
	go func() {
		done <- app(stop)
	}()

	// stoped 判斷當前 stop 狀態
	var stopped bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			log.Printf("server exit err: %+v", err)
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
}

func pprof(stop <-chan struct{}) error {
	// 注意这里主要是为了模拟服务意外退出，用于验证一个服务退出，其他服务同时退出的场景
	go func() {
		server(http.DefaultServeMux, ":8081", stop)
	}()
	time.Sleep(5 * time.Second)
	return fmt.Errorf("mock pprof exit")
}

func app(stop <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	return server(mux, ":8080", stop)
}

// 啟動一個服務
func server(handler http.Handler, addr string, stop <-chan struct{}) error {
	s := http.Server{
		Handler: handler,
		Addr:    addr,
	}

	// 这个 goroutine 我们可以控制退出，因为只要 stop 这个 channel close 或者是写入数据，这里就会退出
	// 同时因为调用了 s.Shutdown 之后，http 这个函数启动的 http server 也会优雅退出
	go func() {
		<-stop
		log.Printf("server will exiting, addr:%s", addr)
		s.Shutdown(context.Background())
	}()

	return s.ListenAndServe()
}

// 紀錄操作異動
// func OSRecord(c *gin.Context, manager string, i interface{}) error {
// 	var act string
// 	inURL := c.FullPath()

// 	op := c.MustGet(admin.TokenAdminInfo).(*admin.Info)

// 	os := &OperationSaveInfo {
// 		OperationItem   string `json:"operation_item"`   // 操作項目 (function.group_id = 0 的  code )
// 		SourceOperation string `json:"source_operation"` // 操作者
// 		TargetOperation string `json:"target_operation"` // 被操作者
// 		Note            string `json:"note"`             // 異動內容
// 		ActType         string `json:"act_type"`
// 		/* 操作動作
// 		"I" // 新增
// 		"U" // 修改
// 		"D" // 刪除
// 		*/
// 	}
// 	os.SourceOperation = op.Account

// 	for _, m := range OperationItemMap() {
// 		if m == manager {
// 			os.OperationItem = m
// 		} else {
// 			return fmt.Errorf("OS Record Should specified a manager")
// 		}
// 	}
// 	// switch manager {
// 	// case OTAM:
// 	// 	os.OperationItem = OTAM
// 	// case OTGM:
// 	// 	os.OperationItem = OTGM
// 	// case OTPM:
// 	// 	os.OperationItem = OTPM
// 	// case OTAgentM:
// 	// 	os.OperationItem = OTAgentM
// 	// case OTCM:
// 	// 	os.OperationItem = OTCM
// 	// case OTLM:
// 	// 	os.OperationItem = OTLM
// 	// case OTSM:
// 	// 	os.OperationItem = OTSM
// 	// case OTCFM:
// 	// 	os.OperationItem = OTCFM
// 	// case OTPRM:
// 	// 	os.OperationItem = OTPRM
// 	// case OTRM:
// 	// 	os.OperationItem = OTRM
// 	// case OTMB:
// 	// 	os.OperationItem = OTMB
// 	// case OTRP:
// 	// 	os.OperationItem = OTRP
// 	// default:
// 	// 	return fmt.Errorf("OS Record Should specified a manager")
// 	// }

// 	switch c.Request.Method {
// 	// case "GET":
// 	// 	act = "READ"
// 	case "POST":
// 		os.ActType = OTActInsert
// 		act = "CREATE"
// 	case "PUT":
// 		os.ActType = OTActUpdate
// 		act = "UPDATE"
// 	case "DELETE":
// 		os.ActType = OTActDelete
// 		act = "DELETE"
// 	}

// 	// 去掉參數
// 	urlStr := strings.Split(inURL, "/")
// 	var r []string
// 	for _, k := range urlStr {
// 		if strings.HasPrefix(k, ":") || strings.HasPrefix(k, "*") {
// 			k = ""
// 		}
// 		if k != "" {
// 			r = append(r, k)
// 		}
// 	}

// 	/*
// 		儲存格式 VERB resource properties note
// 		"Insert: chieftain [manual card info], Note:{id:1565454,name:"name"}"
// 		TODO 統一note格式 by content type
// 	*/
// 	contentType := c.ContentType()
// 	if strings.HasPrefix(contentType, "image/") || strings.HasPrefix(contentType, "audio/") ||
// 		strings.HasPrefix(contentType, "video/") || contentType == "application/octet-stream" {
// 		// if data is binary
// 		os.Note = fmt.Sprintf("%s: %s %s, Note: %s", act, urlStr[1], r[1:], "Data changed")
// 	} else {
// 		os.Note = fmt.Sprintf("%s: %s %s, Note: %s", act, urlStr[1], r[1:], i)
// 	}

// 	s.OperationSave(os)
// 	return nil
// }
