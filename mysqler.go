package main

import (
	"database/sql"
	"fmt"
)

// Initiate Maria Database and defer a Close()
func MariaDBOn() {
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
	// defer dbConnect.Close()
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

// 查詢 Username 並回傳結果
func QueryUser(db *sql.DB, username string) (error, *User) {
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
func DeleteUser(db *sql.DB, username string) error {

	result, err := db.Exec("DELETE FROM users WHERE username=?", username)
	if err != nil {
		fmt.Println("Result: ", result, "Error: ", err)
	}

	return err // usually nil
}

// 檢查是否輸入正確帳號,密碼,二次密碼，並回傳密碼更改是否成功和新密碼
func ChangePassword(db *sql.DB, username, password, password2 string) {
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
