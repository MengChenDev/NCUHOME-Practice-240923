package main

import (
	"Practice-240923/db"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type User struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

func getExistedUsers() []User {
	//检测用户文件是否存在
	_, err := os.Stat(UserFileName)
	if err != nil {
		//创建用户文件
		file, err := os.Create(UserFileName)
		if err != nil {
			fmt.Println("Create user file error:", err)
			return nil
		}
		_, err = file.WriteString("[]")
		if err != nil {
			return nil
		}
		err = file.Close()
		if err != nil {
			return nil
		}
	}

	//读取已有用户
	jsonData, err := os.ReadFile(UserFileName)
	if err != nil {
		fmt.Println("Read user file error:", err)
		return nil
	}

	//解析 JSON 数据
	var users []User
	err = json.Unmarshal(jsonData, &users)
	if err != nil {
		//fmt.Println("Unmarshal user file error:", err)
		users = []User{}
	}

	return users
}

func isExistedUser(account string) (bool, []User) {
	users := getExistedUsers()

	//检查用户是否已存在
	for _, user := range users {
		if user.Account == account {
			return true, users
		}
	}

	return false, users
}

func validateUser(account string, password string) bool {
	users := getExistedUsers()
	for _, user := range users {
		if user.Account == account && user.Password == password {
			return true
		}
	}
	return false
}

func createAccount(account string, password string) bool {
	isExisted, users := isExistedUser(account)

	if isExisted {
		fmt.Println("User already exists")
		return false
	}

	//创建新用户并添加到用户列表
	newUser := User{Account: account, Password: password}
	users = append(users, newUser)

	//将用户列表写入文件
	jsonData, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Marshal user file error:", err)
		return false
	}

	//写入文件
	err = os.WriteFile(UserFileName, jsonData, 0644)
	if err != nil {
		fmt.Println("Write user file error:", err)
		return false
	}

	return true
}

func register() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter account: ")
	accountStr, _ := reader.ReadString('\n')
	account := strings.TrimSuffix(accountStr, "\n")
	fmt.Print("Enter password: ")
	passwordStr, _ := reader.ReadString('\n')
	password := strings.TrimSuffix(passwordStr, "\n")
	if createAccount(account, password) {
		fmt.Println("Account created successfully")
		// 创建文件./dist/account.json并输入{}
		makeDatabaseDir()
		_, err := os.Create(DatabasePathName + "/" + account + ".json")
		if err != nil {
			fmt.Println("Failed to create database file")
			return
		}
		_ = os.WriteFile("./dist/"+account+".json", []byte("{}"), 0644)
	} else {
		fmt.Println("Failed to create account")
	}
}

func login() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please enter your account: ")
	accountStr, _ := reader.ReadString('\n')
	account := strings.TrimSuffix(accountStr, "\n")
	fmt.Print("Please enter your password: ")
	passwordStr, _ := reader.ReadString('\n')
	password := strings.TrimSuffix(passwordStr, "\n")
	if validateUser(account, password) {
		fmt.Println("Login successfully")
		IsLoggedIn = true
		CurrentUser = account
		db.Init(DatabasePathName, account)
		return true
	} else {
		fmt.Println("Account or password is incorrect, or account not registered")
		return false
	}
}
