package main

import (
	"context"
	"log"
	"time"

	pb "acc_service/proto"

	"google.golang.org/grpc"
)

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

// 新增user資料 done
func grpcInsertUser(username, password string) error {
	client := grpcClientOn()
	pbService := pb.NewHelloServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 執行 InsertUser 方法
	r, err := pbService.InsertUser(ctx, &pb.InsertRequest{Username: username,Password:password})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}else{
		log.Printf(r.GetReply())
	}

	grpcClientOff(client)
	return nil	
}

// 查詢 Username 並回傳結果 done
func grpcQueryUser( username string) error {
	client := grpcClientOn()
	pbService := pb.NewHelloServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 執行 QueryUser 方法
	r, err := pbService.QueryUser(ctx, &pb.QueryRequest{Username: username})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}else{
		log.Printf("User found: %s", r.GetReply())
	}

	grpcClientOff(client)
	return nil	
}

// 給予帳號密碼 刪除並回傳成功與否 (會刪除所有相同帳號密碼的人!)
func grpcDeleteUser(username,password string) error {
	client := grpcClientOn()
	pbService := pb.NewHelloServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 執行 QueryUser 方法
	r, err := pbService.DeleteUser(ctx, &pb.DeleteRequest{Username: username,Password:password})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}else{
		log.Printf("User delete: %s", r.GetReply())
	}

	grpcClientOff(client)
	return nil
}

// 檢查是否輸入正確帳號,密碼 並回傳密碼更改是否成功和新密碼,password2 是新密碼
func grpcChangePassword(username, password, newPassword string) error{
	client := grpcClientOn()
	pbService := pb.NewHelloServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 執行 QueryUser 方法
	r, err := pbService.ChangePassword(ctx, &pb.ChangePasswordRequest{Username: username,Password: password,NewPassword: newPassword})
	if err != nil {
		log.Fatalln(err.Error())
		return err
	}else{
		log.Printf(r.GetReply())
	}

	grpcClientOff(client)
	return nil
}
