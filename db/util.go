package db

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	PATH   string
	USER   string
	INITED bool
	DB     map[string]any
)

func Init(path string, user string) {
	PATH = path
	USER = user
	DB = make(map[string]any)
	INITED = true
	readDBFile()
}

// 读取用户数据库文件
func readDBFile() {
	filePath := fmt.Sprintf("./%s/%s.json", PATH, USER)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("User database does not exist", err)
		return
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Failed to read database", err)
		return
	}

	// 解析 JSON 数据
	err = json.Unmarshal(data, &DB)
	if err != nil {
		fmt.Println("Failed to parse database file", err)
		return
	}

	fmt.Println("Connected to database successfully")
}

// 保存用户数据库文件
func saveDBFile() bool {
	filePath := fmt.Sprintf("./%s/%s.json", PATH, USER)

	// 将 DB 转换为 JSON 格式
	data, err := json.Marshal(DB)
	if err != nil {
		fmt.Println("Failed to marshal database to JSON", err)
		return false
	}

	// 写入文件
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		fmt.Println("Failed to write database file", err)
		return false
	}

	return true
}

// Set SET 设置键的值
func Set(key string, value string) bool {
	if !INITED {
		fmt.Println("Database not initialized")
		return false
	}

	// 尝试将 value 解析为 JSON
	var jsonValue interface{}
	err := json.Unmarshal([]byte(value), &jsonValue)
	if err != nil {
		// 如果解析失败，将 value 视为普通字符串
		jsonValue = value
	}

	DB[key] = jsonValue
	fmt.Println("Set key: " + key + " successfully")

	// 保存数据库文件
	if saveDBFile() {
		return true
	}
	return false
}

// Get GET 获取指定的键的值
func Get(key string) (any, bool) {
	if !INITED {
		fmt.Println("Database not initialized")
		return "", false
	}

	value, exists := DB[key]
	if !exists {
		fmt.Println("Key not found: " + key)
		return "", false
	}

	fmt.Println(fmt.Sprintf("%s: %v", key, value))
	return value, true
}

// Del DEL 删除指定的键值对
func Del(key string) bool {
	if !INITED {
		fmt.Println("Database not initialized")
		return false
	}

	_, exists := DB[key]
	if !exists {
		fmt.Println("Key not found: " + key)
		return false
	}

	delete(DB, key)
	fmt.Println("Deleted key: " + key + " successfully")

	if saveDBFile() {
		return true
	}
	return false
}

// SetNX SETNX 如果键不存在，则设置键的值
func SetNX(key string, value string) bool {
	if !INITED {
		fmt.Println("Database not initialized")
		return false
	}

	// 检查键是否已经存在
	_, exists := DB[key]
	if exists {
		fmt.Println("Key already exists: " + key)
		return false
	}

	// 尝试将 value 解析为 JSON
	var jsonValue interface{}
	err := json.Unmarshal([]byte(value), &jsonValue)
	if err != nil {
		// 如果解析失败，将 value 视为普通字符串
		jsonValue = value
	}

	DB[key] = jsonValue
	fmt.Println("Set key: " + key + " successfully")

	if saveDBFile() {
		return true
	}
	return false
}

// LPush LPUSH 将一个或多个值插入到列表的头部
func LPush(key string, values ...string) bool {
	if !INITED {
		fmt.Println("Database not initialized")
		return false
	}

	// 获取当前列表，如果不存在则创建一个新的列表
	list, exists := DB[key].([]interface{})
	if !exists {
		list = make([]interface{}, 0)
	}

	// 将新值添加到列表的尾部
	for _, value := range values {
		var jsonValue interface{}
		err := json.Unmarshal([]byte(value), &jsonValue)
		if err != nil {
			jsonValue = value
		}
		list = append(list, jsonValue)
	}

	DB[key] = list
	fmt.Println("LPUSH to key: " + key + " successfully")

	if saveDBFile() {
		return true
	}
	return false
}

// LRange LRANGE 获取列表中指定范围内的元素
func LRange(key, startStr, stopStr string) ([]interface{}, bool) {
	if !INITED {
		fmt.Println("Database not initialized")
		return nil, false
	}

	start, err := strconv.Atoi(startStr)
	if err != nil {
		fmt.Println("Invalid start value:", err)
		return nil, false
	}

	stop, err := strconv.Atoi(stopStr)
	if err != nil {
		fmt.Println("Invalid stop value:", err)
		return nil, false
	}

	// 获取当前列表
	list, exists := DB[key].([]interface{})
	if !exists {
		fmt.Println("Key not found: " + key)
		return nil, false
	}

	// 计算实际的索引，防止超出范围
	length := len(list)
	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}
	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}

	// 获取范围内的元素
	if start > stop {
		return []interface{}{}, true
	}

	result := list[start : stop+1]
	fmt.Println(fmt.Sprintf("%v", result))
	return result, true
}
