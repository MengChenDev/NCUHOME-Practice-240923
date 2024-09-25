package main

import (
	"Practice-240923/db"
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	DatabasePathName = "dist"
	UserFileName     = "user.json"
)

var (
	IsLoggedIn  = false
	CurrentUser = ""
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter help for more information")
	fmt.Println("Now you can enter some command...")

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "help":
			fmt.Println("help - show this help")
			fmt.Println("exit - exit the program")
			fmt.Println("language - switch language")
			fmt.Println("register - register a new account")
			fmt.Println("login - login to your account")
			fmt.Println("---------------------")
		case "exit":
			break
		case "language":
			fmt.Println("I believe your English proficiency is sufficient to use this command line, so I'm too lazy to develop a Chinese version. After all, who needs a translation when you can just Google it, right?")
		case "register":
			register()
		case "login":
			login()
		default:
			// 处理其他情况
			executeCommand(parts)
		}

	}
}

func executeCommand(parts []string) {
	length := len(parts)
	switch parts[0] {
	case "SET":
		{
			if !checkDatabase() {
				return
			}
			if length != 3 {
				fmt.Println("Invalid command format. Usage: SET <key> <value>")
				return
			}
			db.Set(parts[1], parts[2])
		}
	case "DEL":
		{
			if !checkDatabase() {
				return
			}
			if length != 2 {
				fmt.Println("Invalid command format. Usage: DEL <key>")
				return
			}
			db.Del(parts[1])
		}
	case "SETNX":
		{
			if !checkDatabase() {
				return
			}
			if length != 3 {
				fmt.Println("Invalid command format. Usage: SETNX <key> <value>")
				return
			}
			db.SetNX(parts[1], parts[2])
		}
	case "GET":
		{
			if !checkDatabase() {
				return
			}
			if length != 2 {
				fmt.Println("Invalid command format. Usage: GET <key>")
				return
			}
			db.Get(parts[1])
		}
	case "LPUSH":
		{
			if !checkDatabase() {
				return
			}
			if length < 3 {
				fmt.Println("Invalid command format. Usage: LPUSH <key> <value1> <value2> ...")
				return
			}
			db.LPush(parts[1], parts[2:]...)
		}
	case "LRANGE":
		{
			if !checkDatabase() {
				return
			}
			if length != 4 {
				fmt.Println("Invalid command format. Usage: LRANGE <key> <start> <stop>")
				return
			}
			db.LRange(parts[1], parts[2], parts[3])
		}
	default:
		fmt.Println("Unknown command:", parts[0])
	}
}

func checkDatabase() bool {
	if isInitialized() {
		return true
	} else {
		fmt.Println("Please login first, or register an account")
		fmt.Println("Enter register or login to continue")
		return false
	}
}

func makeDatabaseDir() bool {
	_, err := os.Stat(DatabasePathName)
	if os.IsNotExist(err) {
		err = os.Mkdir(DatabasePathName, 0755)
		if err != nil {
			fmt.Println("Error creating directory:", err)
			return false
		}
	}
	return true
}

func isInitialized() bool {
	if !makeDatabaseDir() {
		return false
	}
	return IsLoggedIn
}
