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

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
)

// db const
var DB *sql.DB
var RC *redis.Client



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
	// router := gin.Default() // 啟動server
	// port := 8088            // port number setting
	// url := ginSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port))
	// router.LoadHTMLGlob("template/html/*")
	// router.Static("/assets", "./template/assets")

	// router.GET("/", GoogleLogin)
	// router.Group("/demo/v1").GET("/hello/:user", hello)
	// router.GET("/login", loginPage)
	// router.POST("/login", loginAuth)
	// router.Group("/oauth").GET("/google/url", LoginWithGoogleOAuth)
	// router.Group("/oauth").GET("/google/login", GoogleLogin)

	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	// router.Run(fmt.Sprintf(":%d", port))

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
