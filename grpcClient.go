package main

import (
	"context"
	"fmt"
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



/////////// gRPC呼叫 gRPC server 進行資料庫操作 ///////////

// 新增user資料 done
func grpcInsertUser(pbService pb.HelloServiceClient,username, password string) error {
	
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

	return nil	
}

// 查詢 Username 並回傳結果 done
func grpcQueryUser(pbService pb.HelloServiceClient, username string) error {
	// client := grpcClientOn()
	// pbService := pb.NewHelloServiceClient(client)
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

	
	return nil	
}

// 給予帳號密碼 刪除並回傳成功與否 (會刪除所有相同帳號密碼的人!)
func grpcDeleteUser(pbService pb.HelloServiceClient,username,password string) error {
	// client := grpcClientOn()
	// pbService := pb.NewHelloServiceClient(client)
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

	
	return nil
}

// 檢查是否輸入正確帳號,密碼 並回傳密碼更改是否成功和新密碼,password2 是新密碼
func grpcChangePassword(pbService pb.HelloServiceClient,username, password, newPassword string) error{
	// client := grpcClientOn()
	// pbService := pb.NewHelloServiceClient(client)
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

	
	return nil
}



// 輸入數字確認訊息來回次數, service 會自行進行 Fibonacci 
func RouteChat(pbService pb.HelloServiceClient,times int) error{
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 執行 RouteChat 方法
	stream , err := pbService.RouteChat(ctx)
	if err != nil {
		return  err
	}

	f := new(pb.RouteNote) // 暫存
	f.Msg = 0
	stream.Send(f)
	
	for  i := 0; i < times; i++ {
		in, err := stream.Recv()
		if err != nil {
			return  err
		}
		// f.Msg = in.Msg
		f.Msg = f.Msg + in.Msg  
		fmt.Println(f.Msg)
	
		stream.Send(f)
	}
	
	// fmt.Println(stream.Trailer())
	return nil
}
