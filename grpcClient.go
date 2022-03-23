package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

var CONNECT grpc.ClientConn

// Initiate and return gRPC Client 
func grpcClientOn()*grpc.ClientConn {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 連線
	addr := "127.0.0.1:8787"
	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err.Error())
	}

	return conn
}
func grpcClientOff(conn *grpc.ClientConn) {
	conn.Close()
}


/////////// gRPC呼叫 gRPC server 進行資料庫操作 ///////////

// 新增user資料
func grpcInsertUser(username, password string) error {
	c := grpcClientOn()
	// 執行 InsertUser 方法
	r, err := c.InsertUser(ctx, &pb.InsertRequest{username: username,password:password})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}else{
		log.Printf("User Inserted: %s", r.GetUsername())
		grpcClientOff(c)
	}
	return err	
}

// 查詢 Username 並回傳結果
func grpcQueryUser(db *sql.DB, username string) (error, *User) {
	user := new(User)
	row := db.QueryRow("SELECT * FROM users WHERE username=?", username)
	if err := row.Scan(&user.ID, &user.User, &user.Password); err != nil {
		fmt.Printf("查詢使用者失敗，原因：%v\n", err)
		return err, user
	}

	fmt.Printf("查詢使用者 ID: %s | Name: %s | Password: %s ", user.ID, user.User, user.Password)
	return nil, user
}

// User mehtod只會刪除並回傳成功與否(不檢查)
func grpcDeleteUser(db *sql.DB, username string) error {

	result, err := db.Exec("DELETE FROM users WHERE username=?", username)
	if err != nil {
		fmt.Println("Result: ", result, "Error: ", err)
	}

	return err // usually nil
}

// 檢查是否輸入正確帳號,密碼,二次密碼，並回傳密碼更改是否成功和新密碼
func grpcChangePassword(db *sql.DB, username, password, password2 string) {
	var user *User
	row := db.QueryRow("SELECT FROM users WHERE username=?", username)
	if err := row.Scan(&user.ID, &user.User, &user.Password, &user.PasswordAgain); err != nil {
		fmt.Printf("查詢使用者失敗，原因：%v\n", err)
	}
	if username == user.User && password == user.Password && password2 == user.PasswordAgain {
		result, err := db.Exec("UPDATE FROM users SET password=? WHERE username=?", password, username)
		if err != nil {
			fmt.Println("UPDATE error: ", err)
		} else {
			fmt.Println("Result: ", result)
		}
	}
}
